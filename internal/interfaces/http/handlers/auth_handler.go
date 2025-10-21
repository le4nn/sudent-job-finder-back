package handlers

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/albkvv/student-job-finder-back/internal/application/usecases"
)

type AuthHandler struct {
    Service *usecases.AuthService
}

func NewAuthHandler(service *usecases.AuthService) *AuthHandler {
    return &AuthHandler{Service: service}
}

func (h *AuthHandler) RegisterPassword(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Password string `json:"password" binding:"required"`
        Role     string `json:"role"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, token, err := h.Service.RegisterPassword(c.Request.Context(), req.Email, req.Phone, req.Password, req.Role)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "token": token,
        "user": gin.H{
            "id":          user.ID,
            "email":       user.Email,
            "phone":       user.Phone,
            "role":        user.Role,
            "is_verified": user.IsVerified,
        },
        "expiresAt": time.Now().Add(24 * time.Hour).Unix(),
    })
}

func (h *AuthHandler) LoginPassword(c *gin.Context) {
    var req struct {
        Identifier string `json:"identifier" binding:"required"`
        Password   string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, token, err := h.Service.LoginPassword(c.Request.Context(), req.Identifier, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":          user.ID,
            "email":       user.Email,
            "phone":       user.Phone,
            "role":        user.Role,
            "is_verified": user.IsVerified,
        },
        "expiresAt": time.Now().Add(24 * time.Hour).Unix(),
    })
}

func (h *AuthHandler) RequestCode(c *gin.Context) {
    var req struct {
        Phone string `json:"phone" binding:"required"`
        Role  string `json:"role" binding:"required,oneof=student employer"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := h.Service.RequestPhoneCode(c.Request.Context(), req.Phone, req.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"result": "sent"})
}

func (h *AuthHandler) RequestEmailCode(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := h.Service.RequestEmailCode(c.Request.Context(), req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "code sent to email"})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
    var req struct {
        Phone string `json:"phone" binding:"required"`
        Code  string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, token, err := h.Service.VerifyPhoneCode(c.Request.Context(), req.Phone, req.Code)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":          user.ID,
            "phone":       user.Phone,
            "role":        user.Role,
            "name":        user.Name,
            "is_verified": user.IsVerified,
        },
        "expiresAt": time.Now().Add(24 * time.Hour).Unix(),
    })
}

func (h *AuthHandler) VerifyEmailCode(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required"`
        Code  string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, token, err := h.Service.VerifyEmailCode(c.Request.Context(), req.Email, req.Code)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":          user.ID,
            "email":       user.Email,
            "role":        user.Role,
            "name":        user.Name,
            "is_verified": user.IsVerified,
        },
        "expiresAt": time.Now().Add(24 * time.Hour).Unix(),
    })
}

