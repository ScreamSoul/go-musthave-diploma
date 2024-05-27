package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/screamsoul/go-musthave-diploma/internal/handlers"
	"github.com/screamsoul/go-musthave-diploma/internal/middlewares"
)

func NewUserLoyaltyRouter(
	uServer *handlers.UserLoyaltyServer,
	globalMiddlewares ...func(http.Handler) http.Handler,
) chi.Router {

	r := chi.NewRouter()

	r.Use(globalMiddlewares...)

	r.Get("/api/ping", uServer.PingStorage)
	r.Post("/api/user/registration", uServer.RegistrationHandler)
	r.Post("/api/user/login", uServer.LoginHandler)

	r.With(middlewares.LoginRequiredMiddleware).Post("/api/user/orders", uServer.UploadOrderHandler)
	r.With(middlewares.LoginRequiredMiddleware).Get("/api/user/orders", uServer.ListOrdersHandler)
	r.With(middlewares.LoginRequiredMiddleware).Get("/api/user/balance", uServer.LoyaltyBalance)
	r.With(middlewares.LoginRequiredMiddleware).Post("/api/user/balance/withdraw", uServer.WithdrawWallet)

	return r
}
