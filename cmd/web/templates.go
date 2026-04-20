package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.ebenezerao.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more /
// to it as the build progresses.

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {

	// Initialize an empty map to act as the cache
	cache := map[string]*template.Template{}

	// filepath.Glob() returns a slice of all filepaths matching the pattern
	// e.g. ["./ui/html/pages/home.tmpl", "./ui/html/pages/view.tmpl"]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through each page template
	for _, page := range pages {

		// Extract just the filename e.g. "home.tmpl", "view.tmpl"
		name := filepath.Base(page)

        // Parse the base template file (base.tmpl) into a template set. 
        ts, err :=template.ParseFiles("./ui/html/base.tmpl")
        if err != nil {
            return nil, err
        }
          // Call ParseGlob() *on this template set* to add any partials.
        ts, err = ts.ParseGlob("./ui/html/partials/ *.tmpl")
        if err != nil {
            return nil, err
        }

         // Call ParseFiles() *on this template set* to add the  page template.
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }


		// Return the map with the page name (e.g. "home.tmpl") as the key and the template set as the value
		cache[name] = ts
	}

	return cache, nil
}
