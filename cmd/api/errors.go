//Filename:cmd/api/errors.go
package main

import(
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error){
	app.logger.PrintError(err, map[string]string{
		"request_method":r.Method,
		"request_url":r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}){
	//create json response
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log the error
	app.logError(r, err)

	//Prepare msg with the error
	message := "the server encountered a problem and couldn't process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The not found response
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//create message
	message := "Request resource couldn't be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

//method not allowed response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	//create our message
	message := fmt.Sprintf("the %s method is not supported for this resources", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// bad request response
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.errorResponse(w, r, http.StatusBadRequest, err.Error())

}

//Validation errors
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

//Edit Conflict errors
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record because of a edit conflict, please try again"
	app.errorResponse(w, r, http.StatusUnprocessableEntity, message)
}
// Rate limit error
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
// Invalid credentials
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Invalid token
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Unauthorized access
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Users who have not activated their account
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// User does not have the required permission (read/write)
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account does not have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}