package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox.alexedwards.net/internal/models"
	"snippetbox.alexedwards.net/ui"
	"time"
)

type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           any
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var function = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page) // Get just the filename, like "home.tmpl"

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(function).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
