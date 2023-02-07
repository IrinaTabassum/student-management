package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := h.stroage.DeleteStudentByID(sID); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	
	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}
