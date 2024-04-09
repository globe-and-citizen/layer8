package utils

import (
	"fmt"
	"net/http"
	"text/template"
)

func ParseHTML(
	w http.ResponseWriter,
	htmlFile string,
	params map[string]interface{},
) {
	t, err := template.New(htmlFile).Delims("[[", "]]").ParseFiles(fmt.Sprintf("assets-v1/templates/%s", htmlFile))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
