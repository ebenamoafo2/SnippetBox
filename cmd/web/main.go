package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide
//  dependencies for the web application.
type application struct {
	logger *slog.Logger
}

 func main() {   

//This is to avoid hardcoding the network address in our code.
 //  Instead, we use the flag package to read the address from a command-line flag.
 addr := flag.String("addr", ":4000", "HTTP network address")
 flag.Parse()

// Use the slog.New() function to initialize a new structured logger, which    
// writes to the standard out stream and uses the default settings.
 logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// Initialize a new instance of our application struct, containing the   
// dependencies (for now, just the structured logger).
 app := &application{
	logger: logger,
 }

 

 



 // Print a log message to say that the server is starting.    
logger.Info("starting server...", "addr", *addr)    

// Use the http.ListenAndServe() function to start a new web server. We pass in    
// two parameters: the TCP network address to listen on (in this case ":4000")   
// and the servemux we just created. If http.ListenAndServe() returns an error    
// we use the log.Fatal() function to log the error message and exit. Note  
// that any error returned by http.ListenAndServe() is always non-nil.   
err := http.ListenAndServe(*addr, app.routes())    
logger.Error(err.Error())
os.Exit(1)
}

//The TCP network address that you pass to http.ListenAndServe() should be in the format "host:port"
//If you omit the host (like we did with ":4000") then the server will listen on all
// your computer’s available network interfaces. Generally, you only need to specify a host in the address if your computer has multiple network interfaces and you want to listen on just one of them.