package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"github.com/juliflorezg/dev-jobs/internal/models"
	"github.com/juliflorezg/dev-jobs/ui"
)

type JobPostsFilterData struct{
	NoPostsData string
	SearchResultMessage []string
	IsSearchResultPage bool
}
type templateData struct {
	CurrentYear int
	JobPosts    []models.JobPost
	JobPostsFilterData
}

func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// that date must be used (https://pkg.go.dev/time@go1.21.6#Time.Format)
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//Extract the file name (like 'home.tmpl.html') and from the full filepath
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
