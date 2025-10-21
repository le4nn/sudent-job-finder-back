package utils

import (
	"fmt"
	"log"
)

func SendSMS(phone, code string) error {
	log.Printf("📱 SMS: Sending OTP code to %s: %s", phone, code)
	fmt.Printf("===================================\n")
	fmt.Printf("SMS OTP CODE\n")
	fmt.Printf("To: %s\n", phone)
	fmt.Printf("Code: %s\n", code)
	fmt.Printf("===================================\n")
	return nil
}
