// Filename:cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	//Create new httprouter instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/forum", app.requirePermission("forums::write", app.createForumHandler))
	router.HandlerFunc(http.MethodGet, "/v1/forum", app.requirePermission("forums:read", app.listForumHandler))
	router.HandlerFunc(http.MethodGet, "/v1/forum/:id", app.requirePermission("forums:read", app.showForumHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/forum/:id", app.requirePermission("forums::write", app.updateForumHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/forum/:id", app.requirePermission("forums::write", app.deleteForumHandler))
	router.HandlerFunc(http.MethodPost, "/v1/comment", app.requirePermission("forums::write", app.createCommentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/comment", app.requirePermission("forums:read", app.listCommentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/comment/:id", app.requirePermission("forums:read", app.showCommentHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/comment/:id", app.requirePermission("forums::write", app.updateCommentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/comment/:id", app.requirePermission("forums::write", app.deleteCommentHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
