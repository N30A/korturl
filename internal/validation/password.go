package validation

// https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
// https://pages.nist.gov/800-63-4/sp800-63b.html#passwordver

import (
	"errors"
	"unicode/utf8"
)

const (
	PasswordMinLength = 15
	PasswordMaxLength = 128
)

var (
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")
)

func Password(password string) error {
	length := utf8.RuneCountInString(password)

	if length < PasswordMinLength {
		return ErrPasswordTooShort
	}

	if length > PasswordMaxLength {
		return ErrPasswordTooLong
	}

	return nil
}
