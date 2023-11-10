package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mfonism/snippetbox/internal/models"
)

type templateData struct {
	IsAuthenticated bool
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
