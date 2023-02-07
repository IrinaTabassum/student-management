package handler

import (
	"log"
	"net/http"
	"student-management/storage"
	"text/template"
)

type StudentForm struct {
	Student   storage.Student
	FormError map[string]error 
	CSRFToken string           
}

func (h Handler) Listofstudent(w http.ResponseWriter, r *http.Request) {

	listStudent, err := h.stroage.ListofStudent()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	t, err := template.ParseFiles("templats/admin/listof-students.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	t.Execute(w, listStudent)
}
