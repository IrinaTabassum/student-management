package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) CreateStudet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create\n\n")
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
			for key, val := range vErr {
				student.FormError[strings.Title(key)] = val
			}

		}
		pareseCreateStudentTemplate(w, student)
		return
	}
	const insertSQuery = `
	INSERT INTO students(
		first_name,
		last_name,
		username,
		email,
		password
	) VALUES (
		:first_name,
		:last_name,
		:username,
		:email,
		:password
	) RETURNING id;	
	`
	stmt, err := h.db.PrepareNamed(insertSQuery)
	if err != nil {
		log.Fatal(err)
	}
	var sid int
	err = stmt.Get(&sid, student)
	if err != nil {
		log.Fatal(err)
	}
	if sid == 0 {
		log.Fatalln("unable to insert student")
	}

	http.Redirect(w, r, "/student/create", http.StatusSeeOther)
}

func pareseCreateStudentTemplate(w http.ResponseWriter, data any) {
	t := template.New("create Student")
	t = template.Must(t.ParseFiles("templats/admin/create-student.html", "templats/admin/_form.html"))
	if err := t.ExecuteTemplate(w, "create-student.html", data); err != nil {
		log.Fatal(err)
	}

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
