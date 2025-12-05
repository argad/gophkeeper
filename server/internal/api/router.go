package api

import (
	"gophkeeper/server/internal/auth"
	"net/http"
)

// NewRouter creates a new router with the given API handlers and JWT manager.
func NewRouter(api *API, jwtManager *auth.JWTManager) http.Handler {
	mux := http.NewServeMux()

	// Public endpoints (no authentication required)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/api/user/register", api.Register)
	mux.HandleFunc("/api/user/login", api.Login)

	// Protected endpoints (authentication required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /api/secrets", api.CreateSecret)
	protectedMux.HandleFunc("GET /api/secrets", api.GetSecrets)
	protectedMux.HandleFunc("GET /api/secrets/{id}", api.GetSecretByID)
	protectedMux.HandleFunc("PUT /api/secrets/{id}", api.UpdateSecret)
	protectedMux.HandleFunc("DELETE /api/secrets/{id}", api.DeleteSecret)

	// Apply JWT authentication middleware to protected routes
	mux.Handle("/api/secrets/", jwtManager.AuthMiddleware(protectedMux))
	mux.Handle("/api/secrets", jwtManager.AuthMiddleware(protectedMux)) // Handle base path too

	return mux
}
