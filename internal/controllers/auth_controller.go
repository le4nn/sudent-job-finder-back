package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/albkvv/student-job-finder-back/internal/models"
	"github.com/albkvv/student-job-finder-back/internal/utils"
)

var codeStore = struct {
	sync.Mutex
	m map[string]phoneCodeData
}{m: make(map[string]phoneCodeData)}

type phoneCodeData struct {
	Code      string
	ExpiresAt time.Time
}

type AuthController struct {
	Users *mongo.Collection
}

func NewAuthController(users *mongo.Collection) *AuthController {
	return &AuthController{Users: users}
}

func (a *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := a.Users.FindOne(ctx, bson.M{"phone": req.Phone}).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, expiresAt, err := utils.GenerateJWT(user.ID.Hex(), 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	resp := models.AuthSession{
		Token: token,
		User: models.UserInfo{
			ID:    user.ID.Hex(),
			Phone: user.Phone,
			Role:  user.Role,
			Name:  user.Name,
		},
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, resp)
}

func (a *AuthController) RequestCode(c *gin.Context) {
	var req models.PhoneCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expires := time.Now().Add(5 * time.Minute)
	codeStore.Lock()
	codeStore.m[req.Phone] = phoneCodeData{Code: code, ExpiresAt: expires}
	codeStore.Unlock()
	fmt.Printf("Код для телефона %s: %s\n", req.Phone, code)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	count, err := a.Users.CountDocuments(ctx, bson.M{"phone": req.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if count == 0 {
		_, err := a.Users.InsertOne(ctx, bson.M{
			"phone":    req.Phone,
			"role":     "student",
			"name":     "",
			"password": "",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "sent"})
}

func (a *AuthController) VerifyCode(c *gin.Context) {
	var req models.PhoneCodeVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	codeStore.Lock()
	data, ok := codeStore.m[req.Phone]
	codeStore.Unlock()
	if !ok || data.ExpiresAt.Before(time.Now()) || data.Code != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong or expired code"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var user models.User
	if err := a.Users.FindOne(ctx, bson.M{"phone": req.Phone}).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	token, expiresAt, err := utils.GenerateJWT(user.ID.Hex(), 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	resp := models.AuthSession{
		Token: token,
		User: models.UserInfo{
			ID:    user.ID.Hex(),
			Phone: user.Phone,
			Role:  user.Role,
			Name:  user.Name,
		},
		ExpiresAt: expiresAt,
	}
	c.JSON(http.StatusOK, resp)
}
