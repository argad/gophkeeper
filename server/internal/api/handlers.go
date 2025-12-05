package api

import (
	"encoding/json"
	"errors"
	"gophkeeper/server/internal/auth"
	"gophkeeper/server/internal/models"
	"gophkeeper/server/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

// API holds the dependencies for the API handlers.
type API struct {
	store      storage.Store
	jwtManager *auth.JWTManager
}

// New creates a new API structure.
func New(store storage.Store, jwtManager *auth.JWTManager) *API {
	return &API{store: store, jwtManager: jwtManager}
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	createdUser, err := a.store.CreateUser(user)
	if err != nil {
		var userExistsErr storage.ErrUserExists
		if errors.As(err, &userExistsErr) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := a.store.GetUserByLogin(creds.Login)
	if err != nil {
		var userNotFoundErr storage.ErrUserNotFound
		if errors.As(err, &userNotFoundErr) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError) // Generic error for other storage issues
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := a.jwtManager.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *API) CreateSecret(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var secret models.Secret
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	secret.UserID = userID // Ensure secret is for the authenticated user

	createdSecret, err := a.store.CreateSecret(secret)
	if err != nil {
		http.Error(w, "Failed to create secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSecret)
}

func (a *API) GetSecrets(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	secrets, err := a.store.GetSecrets(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve secrets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

func (a *API) GetSecretByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// This is a simplified way to get ID from path. A real app might use a router like chi or gorilla/mux
	// to extract path parameters. For net/http, we'd parse r.URL.Path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 { // e.g., /api/secrets/{id}
		http.Error(w, "Invalid secret ID format", http.StatusBadRequest)
		return
	}
	secretIDStr := pathParts[3] // Assuming /api/secrets/123 -> "123"

	secretID, err := strconv.Atoi(secretIDStr)
	if err != nil {
		http.Error(w, "Invalid secret ID", http.StatusBadRequest)
		return
	}

	secret, err := a.store.GetSecretByID(userID, secretID)
	if err != nil {
		var secretNotFoundErr storage.ErrSecretNotFound
		if errors.As(err, &secretNotFoundErr) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve secret", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func (a *API) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var secret models.Secret
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	secret.UserID = userID // Ensure secret is for the authenticated user

	updatedSecret, err := a.store.UpdateSecret(secret)
	if err != nil {
		var secretNotFoundErr storage.ErrSecretNotFound
		if errors.As(err, &secretNotFoundErr) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update secret", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSecret)
}

func (a *API) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 { // e.g., /api/secrets/{id}
		http.Error(w, "Invalid secret ID format", http.StatusBadRequest)
		return
	}
	secretIDStr := pathParts[3]

	secretID, err := strconv.Atoi(secretIDStr)
	if err != nil {
		http.Error(w, "Invalid secret ID", http.StatusBadRequest)
		return
	}

	err = a.store.DeleteSecret(userID, secretID)
	if err != nil {
		var secretNotFoundErr storage.ErrSecretNotFound
		if errors.As(err, &secretNotFoundErr) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}
