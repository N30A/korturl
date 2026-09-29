package validation

import (
	"errors"
	"unicode/utf8"
)

const (
	URLCodeMinLength  = 3
	URLCodeMaxLength  = 64
	URLCodeValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
)

var (
	urlCodeValidCharsMap = func() map[rune]struct{} {
		m := make(map[rune]struct{}, len(URLCodeValidChars))
		for _, char := range URLCodeValidChars {
			m[char] = struct{}{}
		}
		return m
	}()

	reservedURLCodesMap = map[string]struct{}{
		"setup":        {},
		"login":        {},
		"static":       {},
		"api":          {},
		"shorten":      {},
		"shorten/form": {},
	}

	ErrURLCodeTooShort = errors.New("code too short")
	ErrURLCodeTooLong  = errors.New("code too long")
	ErrURLCodeInvalid  = errors.New("code contains invalid characters")
	ErrURLCodeReserved = errors.New("code is reserved")
)

func URLCode(code string) error {
	length := utf8.RuneCountInString(code)

	if length < URLCodeMinLength {
		return ErrURLCodeTooShort
	}

	if length > URLCodeMaxLength {
		return ErrURLCodeTooLong
	}

	for _, char := range code {
		if _, ok := urlCodeValidCharsMap[char]; !ok {
			return ErrURLCodeInvalid
		}
	}

	return nil
}

func ReservedURLCode(code string) error {
	if _, ok := reservedURLCodesMap[code]; ok {
		return ErrURLCodeReserved
	}
	return nil
}
