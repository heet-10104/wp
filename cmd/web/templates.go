package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	_"context"
)

var templateCache map[string]*template.Template

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./resources/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./resources/html/base.tmpl", "./resources/html/partials/partials.tmpl", page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func render(w http.ResponseWriter, status int, page string, data any) {
	ts, ok := templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		serverError(w, err)
		return
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		serverError(w, err)
		return
	}
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func init() {
	var err error
	templateCache, err = newTemplateCache()
	if err != nil {
		panic(err)
	}
}