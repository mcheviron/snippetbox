package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

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

func (v *Validator) CheckField(exist bool, key, message string) {
	if !exist {
		v.AddFieldError(key, message)
	}
}

// To check for empty strings
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// To check the number of characters in a string
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Returns true if a value is in a list of permitted integers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
