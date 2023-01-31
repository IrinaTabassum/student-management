package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete")
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Delete")

	SL, err := getStudentsList()
	if err != nil {
		log.Fatalf("%v", err)
	}
	var newStudentList StudentsList
	for _, student := range SL.Students {
		if student.ID == sID {
			continue
		}
		newStudentList.Students = append(newStudentList.Students, student)
	}

	if err := writeUsersToFile(&newStudentList); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}
