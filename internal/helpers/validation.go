package helpers

import (
	"strings"
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
