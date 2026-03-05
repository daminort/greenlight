package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type ValidationErrors map[string]string

type Validator struct {
	Errors ValidationErrors
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, message string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) >= limit
}

func MaxChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) <= limit
}

func InList[T comparable](value T, list ...T) bool {
	return slices.Contains(list, value)
}

func Matches(value string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(value)
}

func IsEmail(email string) bool {
	return EmailRX.MatchString(email)
}

func IsUnique[T comparable](list []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range list {
		_, seen := uniqueValues[value]
		if seen {
			return false
		}

		uniqueValues[value] = true
	}

	return true
}

func GreaterThan(value int, limit int) bool {
	return value > limit
}

func LessThan(value int, limit int) bool {
	return value < limit
}

func NotZero(value int) bool {
	return value != 0
}

func NotNil(value any) bool {
	return value != nil
}
