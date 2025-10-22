package utils

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPhone(phone string) bool {
	if phone == "" {
		return false
	}

	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")

	phoneRegex := regexp.MustCompile(`^\+?\d{10,15}$`)
	return phoneRegex.MatchString(cleanPhone)
}

func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
