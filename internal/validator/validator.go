package validator

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
	// "github.com/juliflorezg/dev-jobs/"
)

// wh/ Add a new NonFieldErrors []string field to the struct, which we will use to
// hold any validation errors which are not related to a specific form field.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var WebsiteRegex = regexp.MustCompile("^((https?|ftp|smtp):\\/\\/)?(www.)?[a-z0-9]+\\.[a-z]+(\\/[a-zA-Z0-9#]+\\/?)*$")

// Valid() returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key)
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue() returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars() returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

func IsNoData(position, location, contract string) bool {
	if position == "" && location == "" && contract == "" {
		return true
	}
	return false
}

func (v *Validator) ValidateFormData(isNoData bool) {
	if isNoData {
		v.AddFieldError("formError", "Please select at least 1 criteria for your job-post search")
	}
}

func IsValidPassword(pass string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecialChar bool

	for _, c := range pass {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecialChar = true
		}
	}
	fmt.Println(pass)
	fmt.Printf("upper: %v, lower: %v, number: %v, specialChar: %v", hasUpper, hasLower, hasNumber, hasSpecialChar)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecialChar {
		return false
	}

	return true
}
