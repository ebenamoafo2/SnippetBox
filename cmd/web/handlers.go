package main

import (
	"fmt"
	"net/http"
	"strconv"
)

/*
                ================Route pattern ===========
"GET /" - home    => displays the homepage
"GET /snippet/view/{id}" - snippetView => Displays a specific snippet
"GET /snippet/create" - Snippetcreate => Displays a form for creating a new snippet
"POST /snippet/create" - snipptCreatePost => "Saves a new snippet"
*/

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
//w http.ResponseWriter: Used to send a response back to the browser in only bytes. It's a must in Go
// r *http.Request: Contains information about the incoming request.
func home(w http.ResponseWriter, r *http.Request) { 
	w.Header().Add("Server", "Go") 
	  w.Write([]byte("Hello from Snippetbox")) 

} 

//SnippetView handler function
func SnippetView(w http.ResponseWriter, r *http.Request){

// Extract the value of the id wildcard from the request using r.PathValue()    
// and try to convert it to an integer using the strconv.Atoi() function. If    
// it can't be converted to an integer, or the value is less than 1, we    
// return a 404 page not found response.
	
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
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
func SnippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Displays a form for creating a new snippet...."))
}

//SnippetCreatePost handler function
func SnippetCreatePost(w http.ResponseWriter, r *http.Request){
	//Use the w.WriteHeader() method to customize responses sent
	w.WriteHeader(http.StatusCreated)

	//Then use the w.Write() method to write the response body as normal
	w.Write([]byte("Save a new snippet"))
}
