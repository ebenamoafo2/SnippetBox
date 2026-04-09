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
