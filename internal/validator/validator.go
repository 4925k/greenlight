package validator

import "regexp"

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// validator struct with map of validation errors
type Validator struct {
	Errors map[string]string
}

// New returns a new Validator instance
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valids returns true if there are no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// addError adds an error to the map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// check adds and erorr message to the map only if the validation is false
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true is a value is in the given list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func Mathces(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// uniqueValues returns true if values in given slice are unique
func UniqueValues(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
