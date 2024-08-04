package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type M map[string]any

func Template(name string) http.Handler {
	name = filepath.Join("templates", name+".tmpl")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, M{
			"auth": Auth(r),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
