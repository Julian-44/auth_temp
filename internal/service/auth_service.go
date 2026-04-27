package service

import (
	"context"
	"errors"
	"time"

	"auth-server/internal/domain"
	"auth-server/internal/dto"
	"auth-server/internal/repository"
	"auth-server/internal/utils"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtService *utils.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *utils.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(ctx context.Context, input dto.RegisterInput) (*dto.AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, domain.ErrUserExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, input dto.LoginInput) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidPassword
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidPassword
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newAccess, err := s.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	newRefresh, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
