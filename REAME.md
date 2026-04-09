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
