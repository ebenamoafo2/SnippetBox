package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

/*
                ================Route pattern ===========
"GET /" - home    => displays the homepage
"GET /snippet/view/{id}" - snippetView => Displays a specific snippet
"GET /snippet/create" - Snippetcreate => Displays a form for creating a new snippet
"POST /snippet/create" - snipptCreatePost => "Saves a new snippet"
"GET /static/ - http.FileServer => Serves static assets like CSS and JavaScript files"
*/

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
//w http.ResponseWriter: Used to send a response back to the browser in only bytes. It's a must in Go
// r *http.Request: Contains information about the incoming request.
func (app *application) home(w http.ResponseWriter, r *http.Request) { 
	w.Header().Add("Server", "Go") 
// Use the template.ParseFiles() function to read the template file into a    
// template set. If there's an error, we log the detailed error message, use    
// the http.Error() function to send an Internal Server Error response to the    
// user, and then return from the handler so no subsequent code is executed.

	files := []string{
		"ui/html/base.tmpl",
		"ui/html/partials/nav.tmpl",
		"ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

// Then we use the ExecuteTemplate() method on the template set to write the    
// template content as the response body. The last parameter to ExecuteTemplate()    
// represents any dynamic data that we want to pass in, which for now we'll
	err = ts.ExecuteTemplate(w,"base", nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

} 

//SnippetView handler function
func (app *application) SnippetView(w http.ResponseWriter, r *http.Request){

// Extract the value of the id wildcard from the request using r.PathValue()    
// and try to convert it to an integer using the strconv.Atoi() function. If    
// it can't be converted to an integer, or the value is less than 1, we    
// return a 404 page not found response.
	
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.logger.Error("invalid snippet ID", "method", r.Method, "uri", r.URL.RequestURI())
		http.NotFound(w, r)
		return
	}
 // Use the fmt.Sprintf() function to interpolate the id value with a    
 // message, then write it as the HTTP response.   
 	// msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)  
 	// w.Write([]byte(msg))

	fmt.Fprintf(w, "Display a specific snippet with Id %d...", id)

}

//SnippetCreate handler function
func (app *application) SnippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Displays a form for creating a new snippet...."))
}

//SnippetCreatePost handler function
func (app *application) SnippetCreatePost(w http.ResponseWriter, r *http.Request){
	//Use the w.WriteHeader() method to customize responses sent
	w.WriteHeader(http.StatusCreated)

	//Then use the w.Write() method to write the response body as normal
	w.Write([]byte("Save a new snippet"))
}
