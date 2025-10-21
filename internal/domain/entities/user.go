package entities

import "time"

type User struct {
    ID           string
    Email        string
    Phone        string
    PasswordHash string
    Name         string
    Role         string
    IsVerified   bool
    CreatedAt    time.Time
}

type VerificationCode struct {
    Identifier string // email or phone
    Code       string
    Type       string // "email" or "phone"
    Attempts   int
    ExpiresAt  int64 // unix
}

