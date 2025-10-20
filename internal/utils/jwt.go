package utils

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID string `json:"userId"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID string, ttl time.Duration) (string, time.Time, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "dev_secret_change_me"
    }
    expiresAt := time.Now().Add(ttl)
    claims := JWTClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    s, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", time.Time{}, err
    }
    return s, expiresAt, nil
}

func GenerateNumericCode(n int) string {
    max := int64(1)
    for i := 0; i < n; i++ {
        max *= 10
    }
    randNum, _ := rand.Int(rand.Reader, big.NewInt(max))
    return fmt.Sprintf("%0*d", n, randNum.Int64())
}

func LogSMSToConsole(phone, code string) {
    fmt.Printf("SMS-код для номера %s: %s\n", phone, code)
}


