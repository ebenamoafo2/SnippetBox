package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

/*
                ================Route patern ===========
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
	  w.Write([]byte("Hello from Snippetbox")) 

} 

//SnippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request){

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
 	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)  
 	w.Write([]byte(msg))


}

//SnippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Displays a form for creating a new snippet...."))
}

//SnippetCreatePost handler function
func snippetCreatePost(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Save a new snippet"))
}


 func main() {   
// Use the http.NewServeMux() function to initialize a new servemux, then   
 // register the home function as the handler for the "/" URL pattern.    
 mux := http.NewServeMux()    
 mux.HandleFunc("GET /{$}", home) //The {$} is used to restrict the subtree path
 mux.HandleFunc("GET /snippet/view/{id}", snippetView) 
 mux.HandleFunc("GET /snippet/create", snippetCreate)  
 mux.HandleFunc("POST /snippet/create", snippetCreatePost) 


 // Print a log message to say that the server is starting.    
log.Print("starting server on :4000")    

// Use the http.ListenAndServe() function to start a new web server. We pass in    
// two parameters: the TCP network address to listen on (in this case ":4000")   
// and the servemux we just created. If http.ListenAndServe() returns an error    
// we use the log.Fatal() function to log the error message and exit. Note  
// that any error returned by http.ListenAndServe() is always non-nil.   
err := http.ListenAndServe(":4000", mux)    
log.Fatal(err)
}

//The TCP network address that you pass to http.ListenAndServe() should be in the format "host:port"
//If you omit the host (like we did with ":4000") then the server will listen on all
// your computer’s available network interfaces. Generally, you only need to specify a host in the address if your computer has multiple network interfaces and you want to listen on just one of them.