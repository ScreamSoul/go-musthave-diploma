package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

type tokenService struct {
	jwtSecret  string
	jwtExpired time.Duration
}

var signingMethod = jwt.SigningMethodHS256
var service *tokenService

func Initialize(jwtSecret string, jwtExpired time.Duration) {
	service = &tokenService{jwtSecret: jwtSecret, jwtExpired: jwtExpired}
}

func GetTokenSerivce() *tokenService {
	return service
}

func (a *tokenService) GenerateToken(payload *models.Claims) *jwt.Token {
	timeNow := time.Now()

	payload.ExpiresAt = timeNow.Add(a.jwtExpired).Unix()
	payload.IssuedAt = timeNow.Unix()

	return jwt.NewWithClaims(signingMethod, payload)
}

func (a *tokenService) GenerateToString(payload *models.Claims) (string, error) {
	token := a.GenerateToken(payload)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *tokenService) GenerateToCookie(payload *models.Claims) (http.Cookie, error) {
	token := a.GenerateToken(payload)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))

	if err != nil {
		return http.Cookie{}, err
	}

	// Use the token's expiration time for the cookie
	expirationTime := time.Unix(payload.ExpiresAt, 0)
	return http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	}, nil
}

func (a *tokenService) GetUserID(tokenString string) (uuid.UUID, error) {

	claims := models.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// userIdUUID, ok := .(uuid.UUID)
	// if !ok {
	// 	return uuid.Nil, fmt.Errorf("no 'user_id' claim found or invalid UUID format")
	// }
	return claims.UserID, nil
}
