package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Password string             `bson:"password" json:"password"`
	Phone    string             `bson:"phone" json:"phone"`
	Role     string             `bson:"role" json:"role"`
	Name     string             `bson:"name" json:"name"`
}

type LoginRequest struct {
	Phone    string `json:"email"`
	Password string `json:"password"`
}

type PhoneCodeRequest struct {
    Phone string `json:"phone" binding:"required"`
    Role  string `json:"role" binding:"required,oneof=student employer"`
}

type PhoneCodeVerifyRequest struct {
    Phone string `json:"phone" binding:"required"`
    Code string  `json:"code" binding:"required"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

type AuthSession struct {
	Token     string    `json:"token"`
	User      UserInfo  `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
}
