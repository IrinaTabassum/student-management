package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/justinas/nosurf"
)

type LoginFormError struct {
	Username string
	Password string
	UserPass string
}

type LoginStudent struct {
	Username  string
	Password  string
	FormError LoginFormError
	CSRFToken string
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	pareseloginTemplate(w, LoginStudent{

		CSRFToken: nosurf.Token(r),
	})
}
func (h Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	un := r.PostFormValue("Username")
	pass := r.PostFormValue("Password")
	if un == "" {
		pareseloginTemplate(w, LoginStudent{
			Username: un,
			FormError: LoginFormError{
				Username: "The username is required.",
			},
		})
		return
	}
	if pass == "" {
		pareseloginTemplate(w, LoginStudent{
			Username: un,
			FormError: LoginFormError{
				Password: "The password is required.",
			},
		})
		return
	}
	if un != "admin" {
		pareseloginTemplate(w, LoginStudent{
			Username: un,
			Password: "",
			FormError: LoginFormError{
				UserPass: "Incurrect user name or password",
			},
		})
		return
	}
	if pass != "1" {
		pareseloginTemplate(w, LoginStudent{
			Username: un,
			Password: "",
			FormError: LoginFormError{
				UserPass: "Incurrect user name or password",
			},
		})
		return
	}
	h.sessionManager.Put(r.Context(), "username", un)

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)

}

func pareseloginTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templats/login.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.Execute(w, data)
}
