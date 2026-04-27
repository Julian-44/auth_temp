package service

import (
	"auth-server/internal/dto"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterInput) (*dto.AuthResponse, error)
	Login(ctx context.Context, input dto.LoginInput) (*dto.AuthResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
}
