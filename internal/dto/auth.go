// internal/dto/auth.go
package dto

// RegisterInput — входные данные для регистрации
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginInput — входные данные для логина
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshInput — входные данные для обновления токенов
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse — стандартный ответ с токенами
type AuthResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// MessageResponse — ответ с текстовым сообщением
type MessageResponse struct {
	Message string `json:"message"`
}
