package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"student-management/storage"

	"text/template"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/justinas/nosurf"
)

func (h Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	
	editStudent, err := h.stroage.GetStudentByID(sID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	var form StudentForm 
	form.Student = *editStudent
	
	form.CSRFToken = nosurf.Token(r)
	pareseEditUserTemplate(w, form)

}
func (h Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	editstudent :=storage.Student{ID: sID}
	var form StudentForm 
	if err := h.decoder.Decode(&editstudent, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	
	log.Printf("form: %+v \n", editstudent)
    form.Student = editstudent
	
	if err := editstudent.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			newErrs := make(map[string]error)
			for key, val := range vErr {
				newErrs[strings.Title(key)] = val
			}
			form.FormError = newErrs
		}
		pareseCreateStudentTemplate(w, form)
		return
	}
	
	updateStudent, err := h.stroage.UpdateStudent(editstudent)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	 

	http.Redirect(w, r, fmt.Sprintf("/student/%v/edit", updateStudent.ID) , http.StatusSeeOther)

}
func pareseEditUserTemplate(w http.ResponseWriter, data any) {
	t := template.New("edit Student")
	t = template.Must(t.ParseFiles("templats/admin/edit-student.html", "templats/admin/_form.html"))

	t.ExecuteTemplate(w, "edit-student.html", data)

}
