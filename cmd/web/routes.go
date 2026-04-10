package main

import "net/http"



func (app *application) routes() *http.ServeMux {

	
// Use the http.NewServeMux() function to initialize a new servemux, then   
 // register the home function as the handler for the "/" URL pattern.    
 mux := http.NewServeMux() 


	//Basically for serving the css and js files we need to 
// create a file server and then register it as the handler for all URL paths 
// that start with "/static/". For matching paths, we strip the "/static" prefix before the request reaches the file server. 
// This is done using the http.StripPrefix() function.
fileServer := http.FileServer(http.Dir("./ui/static/"))

//Use the mux.Handle() function to register the file server as the handler for
// all URL paths that start with "/static/". For matching paths, we strip the
// "/static" prefix before the request reaches the file server.

mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

//register the other handlers for the remaining URL patterns.
 mux.HandleFunc("GET /{$}", app.home) //The {$} is used to restrict the subtree path
 mux.HandleFunc("GET /snippet/view/{id}", app.SnippetView) 
 mux.HandleFunc("GET /snippet/create", app.SnippetCreate)  
 mux.HandleFunc("POST /snippet/create", app.SnippetCreatePost) 

 return mux
}