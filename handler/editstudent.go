package handler

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

func (h Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	SL, err := getStudentsList()
	if err != nil {
		log.Fatalf("%v", err)
	}
	var editStudent Student
	for _, student := range SL.Students {
		if student.ID == sID {
			editStudent = student
			break
		}
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

	name := r.FormValue("name")
	studentid := r.FormValue("sid")
	bangtla, err := strconv.Atoi(r.FormValue("bang"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	english, err := strconv.Atoi(r.FormValue("engl"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	math, err := strconv.Atoi(r.FormValue("math"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	type ErrorMessage struct {
		Name string
		Role string
	}
	if studentid == "" {
		pareseCreateStudentTemplate(w, ErrorMessage{Role: "The role field is required."})
		return
	}
	if name == "" {
		pareseCreateStudentTemplate(w, ErrorMessage{Name: "The name field is required."})
		return
	}
	studentsList, err := getStudentsList()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for key, student := range studentsList.Students {
		if student.ID == sID {
			studentsList.Students[key].Name = name
			studentsList.Students[key].StudentId = studentid
			studentsList.Students[key].Subject.Bangla = bangtla
			studentsList.Students[key].Subject.English = english
			studentsList.Students[key].Subject.Math = math
		}
	}

	if err := writeUsersToFile(studentsList); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/student/list", http.StatusSeeOther)

}
func pareseEditUserTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templats/edit-student.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.Execute(w, data)
}
