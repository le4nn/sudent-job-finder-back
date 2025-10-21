package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return errors.New("invalid email length")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePhone(phone string) error {
	phone = strings.TrimSpace(phone)
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return errors.New("invalid phone length")
	}
	
	if !phoneRegex.MatchString(cleaned) {
		return errors.New("invalid phone format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func ValidateRole(role string) error {
	if role != "student" && role != "employer" {
		return errors.New("role must be 'student' or 'employer'")
	}
	return nil
}
