package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID    string
	ExpiresAt time.Time
}

func (s *Service) CreateAuthToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	result, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return result, nil

}

func (s *Service) ParseAuthToken(tokenStr string) (*TokenClaims, error) {
	jwtToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}
	ext, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("get expiration time: %w", err)
	}
	return &TokenClaims{
		UserID:    sub,
		ExpiresAt: ext.Time,
	}, nil
}
