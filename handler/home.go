package handler

import (
	"log"
	"net/http"
	"text/template"
)

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templats/index.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.Execute(w, nil)
}
