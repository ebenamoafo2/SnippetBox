# Manipulating the Header Map

In this chapter we used `w.Header().Add()` to add a new header to the response header map. You can also use `Set()`, `Del()`, `Get()`, and `Values()` to manipulate and read values from the header map.

## Example

```go
// Set a new Cache-Control header. If an existing "Cache-Control" header exists,
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")

// In contrast, the Add() method appends a new Cache-Control header and can
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")

// Delete all values for the Cache-Control header.
w.Header().Del("Cache-Control")

// Retrieve the first value for the Cache-Control header.
value := w.Header().Get("Cache-Control")

// Retrieve a slice of all values for the Cache-Control header.
values := w.Header().Values("Cache-Control")
```

## Notes

- `Set()` overwrites existing header values.
- `Add()` appends additional header values.
- `Del()` removes the header entirely.
- `Get()` returns the first header value.
- `Values()` returns all header values as a slice.

## HTML TEMPLATING - CHECK THE HOME.TMPL.THML

Now that we’ve created a template file containing the HTML markup for the home page, the next question is how do we get our home handler to render it?

For this we need to use Go’s `html/template package`, which provides a family of functions for safely parsing and rendering HTML templates. We can use the functions in this package to parse the template file and then execute the template.

## The block action

In the code above we’ve used the {{template}} action to invoke one template from
another. But Go also provides a {{block}}...{{end}} action which you can use instead.
This acts like the {{template}} action, except it allows you to specify some default content
if the template being invoked doesn’t exist in the current template set.
In the context of a web application, this is useful when you want to provide some default
content (such as a sidebar) which individual pages can override on a case-by-case basis if
they need to.

Syntactically you use it like this:
{{define "base"}}

<h1>An example template</h1>
{{block "sidebar" .}}
<p>My default sidebar content</p>
{{end}}
{{end}}

# Serving Static Files in Go

## Overview

In a Go web application, HTML templates handle the structure and content of your pages, but they cannot serve static assets like CSS stylesheets, JavaScript files, or images. These files must be explicitly served using Go's built-in static file serving mechanism.

Go's standard library provides everything needed through three key components in the `net/http` package: `http.FileServer`, `http.Dir`, and `http.StripPrefix`.

---

## Key Components

**http.FileServer** creates an HTTP handler that serves files from a given directory. When it receives a request, it takes the URL path, maps it to a file on disk, and sends the file contents back to the client with the correct Content-Type header.

**http.Dir** converts a string path into a filesystem directory that `http.FileServer` can use. It restricts file serving to the specified directory and its subdirectories.

**http.StripPrefix** wraps another handler and removes a specified prefix from the request URL path before passing the request along. It is essential when the URL prefix used to route requests does not match the directory structure on disk.

---

## How It Works

The key thing to understand is why `StripPrefix` is necessary.

Say your HTML template has this link tag:

```html
<link rel="stylesheet" href="/static/css/styles.css" />
```

The browser sends a request to your server for `/static/css/styles.css`.

**Without StripPrefix**, the FileServer takes the full URL path and appends it to the base directory you provided:

```
./ui/static/ + /static/css/styles.css = ./ui/static/static/css/styles.css  ❌
```

The word "static" appears twice. The file does not exist and the user gets a 404.

**With StripPrefix**, the prefix is removed from the URL before it reaches the FileServer:

```
/static/css/styles.css  →  strips "/static"  →  /css/styles.css
./ui/static/ + /css/styles.css = ./ui/static/css/styles.css  ✅
```

The file is found and served correctly.

---

## Complete Example

**Project structure:**

```
project/
├── main.go
├── ui/
│   ├── static/
│   │   ├── css/
│   │   │   └── styles.css
│   │   └── js/
│   │       └── app.js
│   └── templates/
│       └── index.html
```

**main.go:**

```go
package main

import (
    "html/template"
    "net/http"
)

func main() {
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    http.Handle("/static/", http.StripPrefix("/static", fileServer))

    http.HandleFunc("/", homeHandler)

    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("ui/templates/index.html"))
    tmpl.Execute(w, nil)
}
```

**In your HTML template**, reference static assets using the `/static/` prefix:

```html
<link rel="stylesheet" href="/static/css/styles.css" />
<script src="/static/js/app.js"></script>
```

---

## Official Docs

- `http.FileServer` — https://pkg.go.dev/net/http#FileServer
- `http.StripPrefix` — https://pkg.go.dev/net/http#StripPrefix
- `http.Dir` — https://pkg.go.dev/net/http#Dir

## File server features and functions

Go’s `http.FileServer handler` has a few really nice features that are worth mentioning:
It `sanitizes` all request paths by running them through the `path.Clean()` function before
searching for a file. This removes any . and .. elements from the URL path, which helps
to stop directory traversal attacks. This feature is particularly useful if you’re using the
fileserver in conjunction with a router that doesn’t automatically sanitize URL paths.

`Range requests are fully supported`. This is great if your application is serving large files
and you want to support resumable downloads.

## This is my 6th commit workflow.

## The commit message is :feat(ch3): config flags, structured logging, dependency injection, centralized errors

feat: add configuration, structured logging, dependency injection, and error handling (ch. 3)

This commit covers the full Chapter 3 housekeeping refactor of the Snippetbox app.
No new user-facing features — this is all about making the codebase cleaner,
more configurable, and easier to grow.

---

## 3.1 — Managing configuration settings (cmd/web/main.go)

Replaced the hard-coded server address (":4000") with a command-line flag using
Go's `flag` package. The server address is now configurable at runtime:

    $ go run ./cmd/web -addr=":9999"

Key things to remember:

- `flag.String("addr", ":4000", "HTTP network address")` returns a _pointer_, not a string.
- You MUST call `flag.Parse()` before using any flag value, otherwise you always get the default.
- Dereference the pointer with `*addr` when passing it to `http.ListenAndServe()`.
- `-help` flag is auto-generated and lists all available flags and their defaults — free docs!
- Other flag types: flag.Int(), flag.Bool(), flag.Float64(), flag.Duration().
- Can also use flag.StringVar() to parse directly into a struct field (useful for a config struct).
- env vars alone are worse: no defaults, no type conversion, no -help. Pass them as flags instead:
  $ export SNIPPETBOX_ADDR=":9999" && go run ./cmd/web -addr=$SNIPPETBOX_ADDR

---

## 3.2 — Structured logging (cmd/web/main.go)

Swapped out log.Printf() / log.Fatal() for a custom structured logger using log/slog.

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    logger.Info("starting server", "addr", *addr)
    logger.Error(err.Error())
    os.Exit(1)  // no structured equivalent of log.Fatal(), must do this manually

Log entries now include: timestamp (ms precision), severity level, message, and key-value attributes.

Key things to remember:

- Four severity levels: Debug < Info < Warn < Error. Default minimum level is Info (Debug is silent).
- Attributes are key-value pairs passed after the message: logger.Info("msg", "key", value).
- Mismatched key/value pairs produce !BADKEY in output — use slog.String(), slog.Int() etc. to be safe.
- slog.NewJSONHandler() gives you JSON-formatted logs instead of plaintext.
- Logging to os.Stdout decouples the app from log routing — redirect to a file with >> in prod.
- Custom loggers are concurrency-safe and can be shared across goroutines.
- HandlerOptions lets you set minimum log level or add source file/line info (AddSource: true).

---

## 3.3 — Dependency injection (cmd/web/main.go + handlers.go)

Introduced an `application` struct to hold app-wide dependencies (starting with the logger),
and converted all handler functions into methods on `*application`.

    type application struct {
        logger *slog.Logger
    }

    func (app *application) home(w http.ResponseWriter, r *http.Request) { ... }

Wired up in main():
app := &application{logger: logger}
mux.HandleFunc("GET /{$}", app.home)

Key things to remember:

- Avoids global variables — makes code explicit, testable, and easier to reason about.
- All handlers now access shared dependencies via `app.*` instead of package-level vars.
- This pattern works when all handlers are in the same package.
- If handlers span multiple packages, use a closure pattern with a standalone config package instead.
- This approach scales well — as you add a DB pool, template cache etc., just add fields to the struct.

---

## 3.4 — Centralized error handling (cmd/web/helpers.go)

Created cmd/web/helpers.go with two reusable error helper methods on \*application:

    func (app *application) serverError(w, r, err)  → logs at Error level, sends 500
    func (app *application) clientError(w, status)  → sends specific status code (e.g. 400)

Uses http.StatusText() to produce human-readable status descriptions ("Internal Server Error", etc.).
Updated handlers.go to call app.serverError() instead of inline logging + http.Error().

Key things to remember:

- DRY: no more copy-pasting log + http.Error() in every handler.
- Separation of concerns: error logging logic lives in one place.
- Optional: can include debug.Stack() in serverError() as a "trace" attribute for goroutine stack traces.
- clientError is for user mistakes (bad input, 400, 403, 404), serverError is for our mistakes (500).

---

## 3.5 — Isolating application routes (cmd/web/routes.go)

Extracted all mux/route declarations out of main() into a dedicated routes.go file:

    func (app *application) routes() *http.ServeMux {
        mux := http.NewServeMux()
        // ... all route registrations ...
        return mux
    }

main() now just calls: http.ListenAndServe(\*addr, app.routes())

Key things to remember:

- main() is now focused on exactly 3 things: parse config → wire dependencies → start server.
- routes() is a method on \*application so it can reference app.home, app.snippetView, etc.
- Keeps route declarations in one obvious place — easier to scan and modify later.
- Sets up a good pattern for when routes grow more complex (middleware wrapping, etc.).

---

Files changed:
M cmd/web/main.go — flags, structured logger, application struct, app.routes()
M cmd/web/handlers.go — handlers converted to methods on \*application
A cmd/web/helpers.go — serverError() and clientError() centralized helpers
A cmd/web/routes.go — routes() method extracted from main()
