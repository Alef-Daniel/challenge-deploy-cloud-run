package entity

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidZipcode  = errors.New("invalid zipcode")
	ErrZipcodeNotFound = errors.New("can not find zipcode")
)

var zipcodePattern = regexp.MustCompile(`^\d{8}$`)

func ValidateZipcode(zipcode string) error {
	if !zipcodePattern.MatchString(zipcode) {
		return ErrInvalidZipcode
	}
	return nil
}
