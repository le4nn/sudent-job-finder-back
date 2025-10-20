package entities

type User struct {
    ID    string
    Phone string
    Name  string
    Role  string
}

type VerificationCode struct {
    Phone     string
    Code      string
    ExpiresAt int64 // unix
}

