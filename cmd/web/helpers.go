package main

import (
	"fmt"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page
	//name(like 'home.tmpl'). If no entry exists in the cache with the
	//provided name, then create a new error and call the serverError() helper
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(status)

	//Execute the template set, passing in any dynamic data. If there is an error during execution,
	//  call the serverError() helper.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
