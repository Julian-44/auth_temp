package utils

import (
	"auth-server/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTService(accessSecret, refreshSecret string, accessTTLMinutes, refreshTTLHours int) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     time.Duration(accessTTLMinutes) * time.Minute,
		refreshTTL:    time.Duration(refreshTTLHours) * time.Hour,
	}
}

func (s *JWTService) GenerateAccessToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *JWTService) GenerateRefreshToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

func (s *JWTService) ValidateAccessToken(tokenString string) (int64, error) {
	return s.validateToken(tokenString, s.accessSecret)
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (int64, error) {
	return s.validateToken(tokenString, s.refreshSecret)
}

func (s *JWTService) validateToken(tokenString, secret string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, domain.ErrTokenExpired
		}
		return 0, domain.ErrInvalidToken
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return 0, domain.ErrInvalidToken
}
