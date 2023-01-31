package handler

import (
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
}

func NewHandler(sm *scs.SessionManager, formDecoder *form.Decoder) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formDecoder,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(Method)

	r.Get("/", h.Home)
	r.Get("/login", h.Login)
	r.Post("/login", h.LoginPostHandler)

	r.Route("/student", func(r chi.Router) {
		r.Use(h.Authentication)
		r.Get("/list", h.Listofstudent)
		r.Get("/create", h.CreateStudet)
		r.Post("/store", h.StoreStudent)
		r.Get("/{id:[0-9]+}/edit", h.EditStudent)
		r.Put("/{id:[0-9]+}/update", h.UpdateStudent)
		r.Get("/{id:[0-9]+}/view", h.ViewStudent)
		r.Get("/{id:[0-9]+}/do", h.DeleteStudent)
	})
	r.Group(func(r chi.Router) {
		r.Use(h.Authentication)
		r.Get("/logout", h.LogoutHandler)
	})
	return r
}

func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}
func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		if username == "" {
			// http.Error(w, "unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
