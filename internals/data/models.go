// Filename: internal/data/models.go

package data

import (
	"errors"
	"database/sql"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

// a wrapper for our data models
type Models struct {
	Permissions PermissionModel
	Forums ForumModel
	Comments CommentModel
	Users UserModel
	Tokens TokenModel
}

//NewModels allows us to create a new model
func NewModels(db *sql.DB) Models {
	return Models {
		Forums: ForumModel{DB: db},
		Comments: CommentModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users: UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}