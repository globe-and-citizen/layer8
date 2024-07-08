package utils

import (
	"net/http"
	"path"
	"text/template"
)

func ParseHTML(
	w http.ResponseWriter,
	htmlFile string,
	params map[string]interface{},
) {
	fileName := path.Base(htmlFile)
	t, err := template.New(fileName).Delims("[[", "]]").ParseFiles(htmlFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
