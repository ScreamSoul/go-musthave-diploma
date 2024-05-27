package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/types"
	"go.uber.org/zap"
)

func (s *UserLoyaltyServer) LoyaltyBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	wallet, err := s.loyaltyRepo.GetWalletInfo(r.Context(), userID)
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
