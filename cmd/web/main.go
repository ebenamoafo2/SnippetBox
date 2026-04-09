package main

import (
	"log"
	"net/http"
)


 func main() {   
// Use the http.NewServeMux() function to initialize a new servemux, then   
 // register the home function as the handler for the "/" URL pattern.    
 mux := http.NewServeMux()    
 mux.HandleFunc("GET /{$}", home) //The {$} is used to restrict the subtree path
 mux.HandleFunc("GET /snippet/view/{id}", SnippetView) 
 mux.HandleFunc("GET /snippet/create", SnippetCreate)  
 mux.HandleFunc("POST /snippet/create", SnippetCreatePost) 


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