package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Check whether we have errors accumlated in the correspondent map
func (v *Validator) Empty() bool {
	return len(v.FieldErrors) == 0
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
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
