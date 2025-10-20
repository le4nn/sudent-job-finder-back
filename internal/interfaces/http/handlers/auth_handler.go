package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/albkvv/student-job-finder-back/internal/application/usecases"
)

type AuthHandler struct {
    Service *usecases.AuthService
}

func NewAuthHandler(service *usecases.AuthService) *AuthHandler {
    return &AuthHandler{Service: service}
}

func (h *AuthHandler) RequestCode(c *gin.Context) {
    var req struct { Phone string `json:"phone" binding:"required"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    err := h.Service.RequestCode(req.Phone)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send code"}); return }
    c.JSON(http.StatusOK, gin.H{"result": "sent"})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
    var req struct {
        Phone string `json:"phone" binding:"required"`
        Code  string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    token, user, exp, err := h.Service.VerifyCode(req.Phone, req.Code)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong or expired code"}); return }
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "phone": user.Phone,
            "role": user.Role,
            "name": user.Name,
        },
        "expiresAt": exp,
    })
}

