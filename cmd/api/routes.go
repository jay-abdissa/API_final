//Filename:cmd/api/routes.go
package main

import	(
	"net/http"
	"github.com/julienschmidt/httprouter"
)
func (app *application) routes() http.Handler{
	//Create new httprouter instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed =http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/forum", app.createForumHandler)
	router.HandlerFunc(http.MethodGet, "/v1/forum", app.listForumHandler)
	router.HandlerFunc(http.MethodGet, "/v1/forum/:id", app.showForumHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/forum/:id", app.updateForumHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/forum/:id", app.deleteForumHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comment", app.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comment", app.listCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comment/:id", app.showCommentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comment/:id", app.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comment/:id", app.deleteCommentHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	return app.recoverPanic(app.rateLimit(router))
}