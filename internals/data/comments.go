package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum.castillojadah.net/internals/validator"
)

type Comment struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Content    string    `json:"content"`
	Version    int32     `json:"version"`
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	// Use the Check() method to execute our validation checks

	v.Check(comment.Content != "", "Content", "must be provided")
	v.Check(len(comment.Content) <= 600, "Content", "must not be more than 300 bytes long")
}

// Define a CommentModel which wraps a sql.DB connection pool
type CommentModel struct {
	DB *sql.DB
}

// Insert() allows us  to create a new Comment
func (m CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (content)
		VALUES ($1)
		RETURNING id, created_at, version
	`
	// Collect the data fields into a slice
	args := []interface{}{
		comment.Content,
	}
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt, &comment.Version)
}

// Get() allows us to retrieve a specific Comment
func (m CommentModel) Get(id int64) (*Comment, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, content, version
		FROM comments
		WHERE id = $1
	`
	// Declare a Comment variable to hold the returned data
	var comment Comment

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Content,
		&comment.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &comment, nil
}

// Update() allows us to edit/alter a specific Comment
// Optimistic locking (version number)
func (m CommentModel) Update(comment *Comment) error {
	// Create a query
	query := `
		UPDATE comments
		SET content = $1, version = version + 1
		WHERE id = $2
		AND version = $3
		RETURNING version
	`
	args := []interface{}{
		comment.Content,
		comment.ID,
		comment.Version,
	}

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific Comment
func (m CommentModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM comments
		WHERE id = $1
	`

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// The GetAll() method retuns a list of all the comments sorted by id
func (m CommentModel) GetAll(content string, filters Filters) ([]*Comment, Metadata, error) {

	// Construct the query

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, content, version
		FROM comments
		AND (to_tsvector('simple', content) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortOrder())

	// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{content, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the resultset
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Comment data
	comments := []*Comment{}
	// Iterate over the rows in the resultset
	for rows.Next() {
		var comment Comment
		// Scan the values from the row into comment
		err := rows.Scan(
			&totalRecords,
			&comment.ID,
			&comment.CreatedAt,
			&comment.Content,
			&comment.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Comment to our slice
		comments = append(comments, &comment)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of Comments
	return comments, metadata, nil
}