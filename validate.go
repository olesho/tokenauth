// tokenauth project tokenauth.go
package tokenauth

import (
	"errors"
	"regexp"
)

var mailAddressRegex = regexp.MustCompile(`^[a-zA-Z0-9][-_.a-zA-Z0-9]*@[-_.a-zA-Z0-9]+?$`)
var twoNumbersRegex = regexp.MustCompile(`.*\d.*\d.*`)
var addressRegex = regexp.MustCompile(`.?\d{1,5}.?`)
var phoneRegex = regexp.MustCompile(`^\+?(\s?\d+\s?)+$`)
var fullnameRegex = regexp.MustCompile(`^[A-Za-z]([-']?[A-Za-z]+)*( [A-Za-z]([-']?[A-Za-z]+)*)+$`)

func ValidateEmail(field string) error {
	if !mailAddressRegex.MatchString(field) {
		return errors.New("Invalid email " + field)
	}
	return nil
}

func ValidatePassword(field string) error {
	size := len(field)
	if size < 8 {
		return errors.New("Password cannot be of less then 8 characters size")
	}
	if !twoNumbersRegex.MatchString(field) {
		return errors.New("Password should contain at least two numbers")
	}
	return nil
}

func ValidatePhone(field string) error {
	if !phoneRegex.MatchString(field) {
		return errors.New("Invalid phone number format")
	}
	return nil
}

func ValidateAddress(field string) error {
	if !addressRegex.MatchString(field) {
		return errors.New("Invalid address format")
	}
	return nil
}

func ValidateFullName(field string) error {
	if !fullnameRegex.MatchString(field) {
		return errors.New("Invalid fullname format")
	}
	return nil
}
