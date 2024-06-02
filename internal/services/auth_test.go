package services

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitialize ensures that the Initialize function sets up the service correctly.
func TestInitialize(t *testing.T) {
	Initialize("secret", 1*time.Hour)
	s := GetTokenSerivce()
	assert.NotNil(t, s)
	assert.Equal(t, "secret", s.jwtSecret)
	assert.Equal(t, 1*time.Hour, s.jwtExpired)
}

// TestGenerateToken checks if the token generation works as expected.
func TestGenerateToken(t *testing.T) {
	Initialize("secret", 1*time.Hour)
	s := GetTokenSerivce()

	userID, _ := uuid.NewRandom()
	claims := &models.Claims{
		UserID: userID,
	}

	token := s.GenerateToken(claims)

	require.NotNil(t, token)
	assert.Equal(t, signingMethod.Name, token.Header["alg"])

	// Accessing the custom UserID field directly from our custom Claims struct
	assert.Equal(t, userID, token.Claims.(*models.Claims).UserID)

	// Accessing standard claims directly since they are embedded in our custom Claims struct
	stdClaims := token.Claims.(*models.Claims).StandardClaims
	assert.NotZero(t, stdClaims.ExpiresAt)
	assert.NotZero(t, stdClaims.IssuedAt)
}

func TestGenerateToString(t *testing.T) {
	Initialize("secret", 1*time.Hour)
	s := GetTokenSerivce()

	userID, _ := uuid.NewRandom()
	claims := &models.Claims{
		UserID: userID,
	}

	tokenStr, err := s.GenerateToString(claims)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Optionally, parse the token to verify it contains the correct claims
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, userID.String(), parsedToken.Claims.(jwt.MapClaims)["user_id"])
}

func TestGetUserID(t *testing.T) {
	Initialize("secret", 24*time.Hour)
	service := GetTokenSerivce()

	userID := uuid.New()
	payload := &models.Claims{
		UserID: userID,
	}

	tokenStr, _ := service.GenerateToString(payload)

	parsedUserID, err := service.GetUserID(tokenStr)

	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestGenerateToCookie(t *testing.T) {
	Initialize("secret", 24*time.Hour)
	service := GetTokenSerivce()

	userID := uuid.New()
	payload := &models.Claims{
		UserID: userID,
	}

	cookie, err := service.GenerateToCookie(payload)

	require.NoError(t, err)
	assert.NotNil(t, cookie)
	assert.Equal(t, "token", cookie.Name)
	assert.NotEmpty(t, cookie.Value)
	assert.True(t, cookie.HttpOnly)
	assert.GreaterOrEqual(t, cookie.Expires, time.Now().Add(23*time.Hour))
}
