package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Subject struct {
	Bangla  int `json:"Bangla"`
	English int `json:"English"`
	Math    int `json:"Math"`
}
type Student struct {
	ID        int              `json:"ID"`
	Name      string           `json:"Name"`
	StudentId string           `json:"StudentId"`
	Subject   Subject          `json:"Subject"`
	FormError map[string]error `json:"-" form:"-"`
	CSRFToken string           `json:"-"`
}
type StudentsList struct {
	Students []Student `json:"Students"`
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name field is required"),
		),
	)
}

func (h Handler) Listofstudent(w http.ResponseWriter, r *http.Request) {
	studentList, err := getStudentsList()
	if err != nil {
		log.Fatalf("%v", err)
	}

	t, err := template.ParseFiles("templats/listof-students.html")
	if err != nil {
		log.Fatalf("%v", err)
	}

	t.Execute(w, studentList)
}
func getStudentsList() (*StudentsList, error) {
	f, err := os.Open("data.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var SL StudentsList

	jsonContent, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonContent, &SL); err != nil {
		return nil, err
	}
	return &SL, nil
}
