package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *Creds) UnmarshalJSON(data []byte) error {
	type Alias Creds
	aux := &struct {
		Password string `json:"password"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	hashedPassword := sha256.Sum256([]byte(aux.Password))
	u.Password = hex.EncodeToString(hashedPassword[:])
	return nil
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}
