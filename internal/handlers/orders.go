package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/internal/types"
	"github.com/theplant/luhn"
	"go.uber.org/zap"
)

func (s *UserLoyaltyServer) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		http.Error(w, "Expected Content-Type to be text/plain", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}
	buffer := make([]byte, 50)

	n, err := io.ReadAtLeast(r.Body, buffer, 2)

	if err != nil {
		s.logger.Warn(err.Error(), zap.String("user_id", userID.String()))
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	orderNumver, err := strconv.Atoi(string(buffer[:n]))
	if err != nil || !luhn.Valid(orderNumver) {
		s.logger.Warn("invalid order number format", zap.String("user_id", userID.String()))
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	err = s.store.CreateOrder(r.Context(), orderNumver, userID)

	if errors.Is(err, repositories.ErrOrderAlreadyUpload) {
		w.WriteHeader(http.StatusOK)
	} else if errors.Is(err, repositories.ErrOrderConflict) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		s.logger.Error(err.Error(), zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *UserLoyaltyServer) ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserID).(uuid.UUID)

	if !ok {
		s.logger.Error("Failed to extract user ID from context", zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	orders, err := s.store.ListOrders(r.Context(), userID)

	if err != nil {
		s.logger.Error(err.Error(), zap.String("user_id", userID.String()))
		http.Error(w, "Interanl error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			s.logger.Error("Error writing response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
