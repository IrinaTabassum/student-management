package handler

import (
	"fmt"
	"student-management/storage"


	"log"
	"net/http"
	"strings"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) CreateStudet(w http.ResponseWriter, r *http.Request) {
	pareseCreateStudentTemplate(w, StudentForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	form := StudentForm{}
	student := storage.Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	log.Printf("form: %+v \n", student)
	form.Student=student

	if err := student.Validate(); err != nil {
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
	
	storStudent, err := h.stroage.CreateStudent(student)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/student/%v/edit", storStudent.ID) , http.StatusSeeOther)
}
  
func pareseCreateStudentTemplate(w http.ResponseWriter, data any) {
	t := template.New("create Student")
	t = template.Must(t.ParseFiles("templats/admin/create-student.html", "templats/admin/_form.html"))
	if err := t.ExecuteTemplate(w, "create-student.html", data); err != nil {
		log.Fatal(err)
	}

}
 
