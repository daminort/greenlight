package users

import (
	"greenlight.damian.net/internal/pkg/validator"
)

func validatePassword(v *validator.Validator, pwd string) {
	v.Check(validator.NotBlank(pwd), "password", "must be provided")
	v.Check(validator.MinChars(pwd, 8), "password", "must be at least 8 characters long")
	v.Check(validator.MaxChars(pwd, 72), "password", "must not be more than 72 characters long")
}

func ValidateUser(u *User) *validator.Validator {
	v := validator.New()

	v.Check(validator.NotBlank(u.Name), "name", "must be provided")
	v.Check(validator.MaxChars(u.Name, 500), "name", "must not be more than 500 characters")

	v.Check(validator.NotBlank(u.Email), "email", "must be provided")
	v.Check(validator.IsEmail(u.Email), "email", "must be a valid email")

	if u.Pwd.Text != nil {
		validatePassword(v, *u.Pwd.Text)
	}

	if u.Pwd.Hash == nil {
		panic("missing password hash for user")
	}

	return v
}
