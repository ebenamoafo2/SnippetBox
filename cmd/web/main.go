package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.ebenezerao.net/internal/models"
)

// Define an application struct to hold the application-wide
//
//	dependencies for the web application.
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel //Added to our struct to make it available to all our handlers
}

func main() {

	//This is to avoid hardcoding the network address in our code.
	//  Instead, we use the flag package to read the address from a command-line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// The dsn flag holds the MySQL data source name, which is a string that contains
	// the information needed to connect to a MySQL database. The format of the DSN is "username:password@protocol(address)/dbname?param=value".
	// In this case, we are connecting to a MySQL database with the username "web", password "pass", and database name "snippetbox". The parseTime=true parameter tells the MySQL driver to parse time values into Go's time.Time type.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}

	//run this later, when the function exits ensuring resource leaks are prevented
	defer db.Close()

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	// Print a log message to say that the server is starting.
	logger.Info("starting server...", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close() // Close the database connection if the ping fails to free up resources.
		return nil, err
	}
	return db, nil
}

//The TCP network address that you pass to http.ListenAndServe() should be in the format "host:port"
//If you omit the host (like we did with ":4000") then the server will listen on all
// your computer’s available network interfaces. Generally, you only need to specify a host in the address if your computer has multiple network interfaces and you want to listen on just one of them.
