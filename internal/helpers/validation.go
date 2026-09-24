package helpers

import (
	"regexp"
	"strings"
	"unicode"
	"uuid"
)

func CheckIfStringEmpty(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	return false
}

func CheckIfUuid(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}

	return true
}

func CheckStringLen(s string, minLen int64, maxLen int64) bool {
	length := int64(len(strings.TrimSpace(s)))
	return length >= minLen && length <= maxLen
}

func CheckIfValidNumber(num float32, required bool) bool {
	if required && num <= 0 {
		return false
	}

	return true
}

func CheckArrayEmpty[T any](value []T) bool {
	return len(value) == 0
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// CheckIfValidEmail returns true if the email matches a standard email pattern.
func CheckIfValidEmail(email string) bool {
	if CheckIfStringEmpty(email) {
		return false
	}
	return emailRegex.MatchString(email)
}

// CheckIfStrongPassword returns true if password is at least 8 chars and
// contains at least one uppercase letter, one number, and one special character.
func CheckIfStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	return hasUpper && hasNumber && hasSpecial
}
