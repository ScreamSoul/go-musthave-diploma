package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/internal/types"
	"github.com/theplant/luhn"
	"go.uber.org/zap"
)

func (s *UserLoyaltyServer) LoyaltyBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	wallet, err := s.loyaltyRepo.GetWallet(r.Context(), userID)
	if err != nil {
		s.logger.Error(err.Error(), zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(wallet); err != nil {
		s.logger.Error("Error writing response", zap.Error(err))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserLoyaltyServer) WithdrawWallet(w http.ResponseWriter, r *http.Request) {
	var withdraw models.Withdraw

	contentType := r.Header.Get("Content-Type")
	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	} else if contentType != "application/json" {
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	} else if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	} else if !luhn.Valid(int(withdraw.Order)) {
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	err := s.loyaltyRepo.WithdrawWallet(
		r.Context(),
		userID,
		&withdraw,
	)
	if errors.Is(err, repositories.ErrLowBalance) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}
	if err != nil {
		s.logger.Error(err.Error(), zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserLoyaltyServer) ListWithdraws(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	withdraws, err := s.loyaltyRepo.GetWithdrawals(r.Context(), userID)

	if err != nil {
		s.logger.Error(err.Error(), zap.String("user_id", userID.String()))

		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(withdraws) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(withdraws); err != nil {
		s.logger.Error("Error writing response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
