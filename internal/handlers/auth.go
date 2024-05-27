package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/internal/services"
	"go.uber.org/zap"
)

func (s *UserLoyaltyServer) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Creds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		s.logger.Error(err.Error(), zap.Any("body", creds))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userID, err := s.loyaltyRepo.CreateUser(r.Context(), &creds)

	if err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			s.logger.Error(err.Error(), zap.Any("body", creds))
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	tokenService := services.GetTokenSerivce()

	cookie, err := tokenService.GenerateToCookie(&models.Claims{UserID: userID})
	if err != nil {
		s.logger.Error("generate token error", zap.Error(err))
		http.Error(w, "Interanal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}

func (s *UserLoyaltyServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Creds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := s.loyaltyRepo.CheckUserPassword(r.Context(), &creds)

	if err != nil {
		if errors.Is(err, repositories.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	tokenService := services.GetTokenSerivce()

	cookie, err := tokenService.GenerateToCookie(&models.Claims{UserID: userID})
	if err != nil {
		s.logger.Error("generate token error", zap.Error(err))
		http.Error(w, "Interanal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
