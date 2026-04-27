// internal/delivery/http/handler/auth_handler.go
package handler

import (
	"auth-server/internal/dto"
	"auth-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "invalid input: " + err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		// обработка доменных ошибок
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "invalid input"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
