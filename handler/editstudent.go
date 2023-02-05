package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"text/template"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var editStudent Student
	const editStudentQuery = `SELECT * FROM students WHERE id=$1 AND deleted_at IS NULL`
	if err := h.db.Get(&editStudent, editStudentQuery, id); err != nil {
		log.Fatal(err)
	}

	editStudent.CSRFToken = nosurf.Token(r)
	pareseEditUserTemplate(w, editStudent)

}
func (h Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := r.ParseForm(); err != nil {
		log.Fatalf("%#v", err)
	}

	editstudent := Student{ID: sID}

	if err := h.decoder.Decode(&editstudent, r.PostForm); err != nil {
		log.Fatal(err)
	}
	
	log.Printf("form: %+v \n", editstudent)

	if err := editstudent.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			newErrs := make(map[string]error)
			for key, val := range vErr {
				newErrs[strings.Title(key)] = val
			}
			editstudent.FormError = newErrs
		}
		pareseCreateStudentTemplate(w, editstudent)
		return
	}
	const updatStudenrQuery =`
	UPDATE students SET
		first_name = :first_name,
		last_name = :last_name,
		status = :status
	WHERE id = :id AND deleted_at IS NULL;
	`
	stud, err := h.db.PrepareNamed(updatStudenrQuery)
	if err != nil{
		log.Fatal(err)
	}
	res, err := stud.Exec(editstudent)
	if err != nil{
		log.Fatal(err)
	}
	rocount, err := res.RowsAffected()
	if err != nil{
		log.Fatal(err)
	}

	if rocount >0 {
		http.Redirect(w, r, "/student/list", http.StatusSeeOther)
		return
	}

	pareseCreateStudentTemplate(w, editstudent)

}
func pareseEditUserTemplate(w http.ResponseWriter, data any) {
	t := template.New("edit Student")
	t = template.Must(t.ParseFiles("templats/admin/edit-student.html", "templats/admin/_form.html"))

	t.ExecuteTemplate(w, "edit-student.html", data)

}
