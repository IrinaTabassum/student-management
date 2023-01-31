package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/justinas/nosurf"
)

func (h Handler) CreateStudet(w http.ResponseWriter, r *http.Request) {
	pareseCreateStudentTemplate(w, Student{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("%#v", err)
	}

	student := Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Fatal(err)
	}
	log.Printf("form: %+v \n", student)

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			fmt.Println(vErr)
			student.FormError = vErr
			fmt.Println(student.FormError)
		}
		pareseCreateStudentTemplate(w, student)
		return
	}

	studentsList, err := getStudentsList()
	if err != nil {
		log.Fatalf("%v", err)
	}

	lenOfSL := len(studentsList.Students)
	if lenOfSL > 0 {
		student.ID = studentsList.Students[lenOfSL-1].ID + 1
	} else {
		student.ID = 1
	}

	studentsList.Students = append(studentsList.Students, student)

	if err := writeUsersToFile(studentsList); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/student/create", http.StatusSeeOther)
}

func pareseCreateStudentTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templats/create-student.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.Execute(w, data)
}
func writeUsersToFile(SL *StudentsList) error {
	jsonContent, err := json.MarshalIndent(SL, " ", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data.json", jsonContent, 0644)
	if err != nil {
		return err
	}
	return nil
}
