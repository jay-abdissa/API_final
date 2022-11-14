//Filename:cmd/api/handlers.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"forum.castillojadah.net/internals/data"
	"forum.castillojadah.net/internals/validator"
)

// createCommentHandler for the "POST /v1/comment" endpoint
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Our target decode destination
	var input struct {
		Content  string `json:"content"`
	}
	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Comment struct
	comment := &data.Comment{
		Content:  input.Content,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateComment(v, comment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Create a Comment
	err = app.models.Comments.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Create a Location header for the newly created resource/Comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comment/%d", comment.ID))
	// Write the JSON response with 201 - Created status code with the body
	// being the Comment data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showCommentHandler for the "GET /v1/comment/:id" endpoint
func (app *application) showCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific comment
	comment, err := app.models.Comments.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the comment that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	comment, err := app.models.Comments.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		Content  *string `json:"content"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check for updates
	if input.Content != nil {
		comment.Content = *input.Content
	}
	
	// Perform validation on the updated Comment. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateComment(v, comment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated Comment record to the Update() method
	err = app.models.Comments.Update(comment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the comment that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Comment from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Comments.Delete(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "comment item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listCommentHandler() allows the client to see a listing of comments
// based on a set of criteria
func (app *application) listCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Create an input struct to hold our query parameters
	var input struct {
		Content string
		data.Filters
	}
	// Initialize a validator
	v := validator.New()
	// Get the URL values map
	qs := r.URL.Query()
	// Use the helper methods to extract the values
	input.Content = app.readString(qs, "content", "")
	// Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specific the allowed sort values
	input.Filters.SortList = []string{"id", "content", "-id", "-content"}
	// Check for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Get a listing of all comments
	comments, metadata, err := app.models.Comments.GetAll(input.Content, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containg all the comments
	err = app.writeJSON(w, http.StatusOK, envelope{"comments": comments, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}