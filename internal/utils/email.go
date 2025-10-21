package utils

import (
	"fmt"
	"log"
)

func SendEmail(email, code string) error {
	log.Printf("📧 EMAIL: Sending OTP code to %s: %s", email, code)
	fmt.Printf("===================================\n")
	fmt.Printf("EMAIL OTP CODE\n")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Code: %s\n", code)
	fmt.Printf("===================================\n")
	return nil
}
