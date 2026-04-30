package main

import (
	"fmt"
	"net/http"
)

// commonHeaders adds security-related HTTP headers to every response:
//   - Content-Security-Policy:  Restricts resource loading to same origin, plus Google Fonts
//   - Referrer-Policy:          Only sends origin (not full URL) to cross-origin requests
//   - X-Content-Type-Options:   Stops browsers from MIME-sniffing responses
//   - X-Frame-Options:          Blocks the site from being iframed (prevents clickjacking)
//   - X-XSS-Protection:         Disabled intentionally — the old browser XSS filter can cause its own vulnerabilities
//   - Server:                   Hides the real server info, just reports "Go"


func commonHeaders(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy","default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options","deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w,r)
	})
}

//LogRequest records the IP address of the user, and the method, URI and HTTP version for the request.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip = r.RemoteAddr
			proto = r.Proto
			method = r.Method
			uri = r.URL.RequestURI()

		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w,r)
	})
}

// recoverPanic catches any panics in downstream handlers, closes the connection,
// and returns a 500 error to the client instead of crashing the goroutine.
func (app *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Create a deferred function (which will always be run in the event        
// of a panic as Go unwinds the stack).	
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w,r,fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w,r )
	})
}