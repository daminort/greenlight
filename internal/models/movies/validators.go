package movies

import (
	"time"

	"greenlight.damian.net/internal/validator"
)

func ValidateMovie(m *Movie) *validator.Validator {
	v := validator.New()

	v.Check(validator.NotBlank(m.Title), "title", "must be provided")
	v.Check(validator.MaxChars(m.Title, 50), "title", "must not be more than 50 characters")

	v.Check(validator.NotZero(m.Year), "year", "must be provided")
	v.Check(validator.GreaterThan(m.Year, 1887), "year", "must be greater than or equal to 1888")
	v.Check(validator.LessThan(m.Year, int(time.Now().Year())+1), "year", "must not be in the future")

	v.Check(validator.NotZero(int(m.Runtime)), "runtime", "must be provided")
	v.Check(validator.GreaterThan(int(m.Runtime), 0), "runtime", "must be positive")

	v.Check(validator.NotNil(m.Genres), "genres", "must be provided")
	v.Check(validator.GreaterThan(len(m.Genres), 0), "genres", "must contain at least 1 genre")
	v.Check(validator.LessThan(len(m.Genres), 6), "genres", "must not contain more than 5 genres")
	v.Check(validator.IsUnique(m.Genres), "genres", "must not contain duplicate values")

	return v
}
