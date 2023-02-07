package storage

import (
	"database/sql"
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