package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	
)

type Student struct {
	ID        int              `db:"id"`
	FirstName string           `db:"first_name"`
	LastName  string           `db:"last_name"`
	Email     string           `db:"email"`
	Username  string           `db:"username"`
	Password  string           `db:"password"`
	Status    bool             `db:"status"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
	DeletedAt sql.NullTime     `db:"deleted_at"`
	FormError map[string]error `form:"-"`
	CSRFToken string           `form:"csrf_token"`
}
type StudentsList struct {
	Students []Student `json:"Students"`
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FirstName,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("The last name field is required."),
		),
		validation.Field(&s.Username,
			validation.Required.When(s.ID == 0).Error("The username field is required."),
		),
		validation.Field(&s.Email,
			validation.Required.When(s.ID == 0).Error("The email field is required."),
			is.Email.Error("The email field must be a valid email."),
		),
		validation.Field(&s.Password,
			validation.Required.When(s.ID == 0).Error("The password field is required."),
		),
	)
}

func (h Handler) Listofstudent(w http.ResponseWriter, r *http.Request) {
	const listQuery = `SELECT * from students WHERE deleted_at IS NULL`
	var listStudent []Student

	if err := h.db.Select(&listStudent, listQuery); err != nil {
		log.Fatalln(err)
	}

	t, err := template.ParseFiles("templats/admin/listof-students.html")
	if err != nil {
		log.Fatalf("%v", err)
	}

	t.Execute(w, listStudent)
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
