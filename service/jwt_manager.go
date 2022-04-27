package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTManager is a JSON web token manager
type JWTManager struct {
	secretKey string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token)(interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok{
				return nil, fmt.Errorf("unexpected token signing method")
			}


			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok{
		return nil, fmt.Errorf("Invalid token claims ", err)
	}

	return claims, nil
}