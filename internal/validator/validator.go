package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Check whether we have errors accumlated in the correspondent map
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = message
	}
}

// Check whether a certain error occured and set the correlating fields accordingly
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Check if the provided string is empty
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Check whether the letters in a string exceed a certain value (n)
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Returns true if a value is in a list of permitted integers
// eg. 7 is in a list of (1, 7, 365)
// func PermittedInt(value int, permittedValues ...int) bool {
// 	for i := range permittedValues {
// 		if value == permittedValues[i] {
// 			return true
// 		}
// 	}
// 	return false
// }

// NOTE: replaced PermittedInt() with this generic function for feasibility
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Make sure that the password length isn't less than a certain number
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Check if the email is of a valid format
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
