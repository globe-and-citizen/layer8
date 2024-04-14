package utils

import (
	"net/http"
	"strings"
	"text/template"
)

func ParseHTML(
	w http.ResponseWriter,
	htmlFile string,
	params map[string]interface{},
) {
	t, err := template.New(htmlFile[strings.LastIndex(htmlFile, "/")+1:]).Delims("[[", "]]").ParseFiles(htmlFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
