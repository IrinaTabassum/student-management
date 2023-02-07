package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h Handler) ViewStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Panicln(sID)
	// SL, err := getStudentsList()
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// var viewStudent Student
	// for _, student := range SL.Students {
	// 	if student.ID == sID {
	// 		viewStudent = student
	// 		break
	// 	}
	// }
	// t, err := template.ParseFiles("templats/view-student.html")
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// t.Execute(w, viewStudent)

}
