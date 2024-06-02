package middlewares

import (
	"context"
	"net/http"

	"github.com/screamsoul/go-musthave-diploma/internal/services"
	"github.com/screamsoul/go-musthave-diploma/internal/types"
	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
)

func LoginRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()

		cookie, err := r.Cookie("token")
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		tokenService := services.GetTokenSerivce()
		userID, err := tokenService.GetUserID(cookie.Value)
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, "invalide token", http.StatusUnauthorized)
			return
		}

		// Создаем новый контекст с user_id
		ctx := context.WithValue(r.Context(), types.UserID, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
