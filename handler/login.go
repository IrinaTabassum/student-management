package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"student-management/storage/postgres"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
)

type LoginStudent struct {
	Username  string
	Password  string
	FormError map[string]error
	CSRFToken string
}

func (ls LoginStudent) Validate() error {
	return validation.ValidateStruct(&ls,
		validation.Field(&ls.Username,
			validation.Required.Error("The username field is required."),
		),
		validation.Field(&ls.Password,
			validation.Required.Error("The password field is required."),
		),
	)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	pareseloginTemplate(w, LoginStudent{
		CSRFToken: nosurf.Token(r),
	})
}
func (h Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var loginStudent LoginStudent
	if err := h.decoder.Decode(&loginStudent, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := loginStudent.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			formErr := make(map[string]error)
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
			loginStudent.FormError = formErr
			loginStudent.Password = ""
			loginStudent.CSRFToken = nosurf.Token(r)
			pareseloginTemplate(w, loginStudent)
			return
		}
	}
	getStudent, err := h.stroage.GetStudentByUsername(loginStudent.Username)

	if err != nil {
		if err.Error() == postgres.NotFound {
			formErr := make(map[string]error)
			formErr["Username"] = fmt.Errorf("credentials does not match")
			loginStudent.FormError = formErr
			loginStudent.CSRFToken = nosurf.Token(r)
			loginStudent.Password = ""
			pareseloginTemplate(w, loginStudent)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(getStudent.Password), []byte(loginStudent.Password)); err != nil {
		formErr := make(map[string]error)
		formErr["Username"] = fmt.Errorf("credentials does not match")
		loginStudent.FormError = formErr
		loginStudent.CSRFToken = nosurf.Token(r)
		loginStudent.Password = ""
		pareseloginTemplate(w, loginStudent)
		return
	}

	h.sessionManager.Put(r.Context(), "studentID", strconv.Itoa(getStudent.ID))

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)

}

func pareseloginTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templats/login.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.Execute(w, data)
}
