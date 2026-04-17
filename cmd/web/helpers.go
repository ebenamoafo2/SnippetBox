package main

import (
	"net/http"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	//HTTP method -> r.Method eg. GET, POST, etc.
	//Requested URL -> /home,/snippet/view/1, etc.
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// The clientError helper sends a specific status code and corresponding description
// to the user.

func (app *application) clientError(w http.ResponseWriter, status int) {
	//statusText() function returns a text description for the provided HTTP status code.
	// This returns a human-friendly text representation of a given HTTP status code —
	// for example http.StatusText(400) will return the string "Bad Request", and http.StatusText(500) will return the string "Internal Server Error".
	http.Error(w, http.StatusText(status), status)
}
