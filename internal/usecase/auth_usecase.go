package usecase

import (
	"auth-server/internal/domain"
	"auth-server/internal/repository"
	"auth-server/internal/utils"
	"context"
	"time"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtService *utils.JWTService
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtService *utils.JWTService) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCase) Register(ctx context.Context, email, password string) (int64, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return 0, domain.ErrUserExists
	}

	// Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", "", domain.ErrInvalidPassword
		}
		return "", "", err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", domain.ErrInvalidPassword
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *authUseCase) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Verify user still exists
	_, err = uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := uc.jwtService.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
