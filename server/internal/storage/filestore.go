package storage

import (
	"gophkeeper/server/internal/models"
	"sync"
)

// MemStore is an in-memory data store.
type MemStore struct {
	mu           sync.RWMutex
	users        map[string]models.User  // map[login]User
	secrets      map[int][]models.Secret // map[userID][]Secret
	nextUserID   int
	nextSecretID int
}

// NewMemStore creates and returns a new MemStore.
func NewMemStore() *MemStore {
	return &MemStore{
		users:        make(map[string]models.User),
		secrets:      make(map[int][]models.Secret),
		nextUserID:   1,
		nextSecretID: 1,
	}
}

// CreateUser adds a new user to the store.
func (s *MemStore) CreateUser(user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Login]; exists {
		return models.User{}, NewErrUserExists(user.Login)
	}

	user.ID = s.nextUserID
	s.users[user.Login] = user
	s.nextUserID++
	return user, nil
}

// GetUserByLogin retrieves a user by their login.
func (s *MemStore) GetUserByLogin(login string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[login]
	if !exists {
		return models.User{}, NewErrUserNotFound(login)
	}
	return user, nil
}

// CreateSecret adds a new secret for a user.
func (s *MemStore) CreateSecret(secret models.Secret) (models.Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret.ID = s.nextSecretID
	s.secrets[secret.UserID] = append(s.secrets[secret.UserID], secret)
	s.nextSecretID++
	return secret, nil
}

// GetSecrets retrieves all secrets for a specific user.
func (s *MemStore) GetSecrets(userID int) ([]models.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userSecrets, exists := s.secrets[userID]
	if !exists {
		return []models.Secret{}, nil
	}
	return userSecrets, nil
}

// GetSecretByID retrieves a specific secret for a user by its ID.
func (s *MemStore) GetSecretByID(userID, secretID int) (models.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if userSecrets, exists := s.secrets[userID]; exists {
		for _, secret := range userSecrets {
			if secret.ID == secretID {
				return secret, nil
			}
		}
	}
	return models.Secret{}, NewErrSecretNotFound(secretID)
}

// UpdateSecret updates an existing secret for a user.
func (s *MemStore) UpdateSecret(secret models.Secret) (models.Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userSecrets, exists := s.secrets[secret.UserID]; exists {
		for i, sct := range userSecrets {
			if sct.ID == secret.ID {
				s.secrets[secret.UserID][i] = secret
				return secret, nil
			}
		}
	}
	return models.Secret{}, NewErrSecretNotFound(secret.ID)
}

// DeleteSecret deletes a secret for a user by its ID.
func (s *MemStore) DeleteSecret(userID, secretID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userSecrets, exists := s.secrets[userID]; exists {
		for i, secret := range userSecrets {
			if secret.ID == secretID {
				s.secrets[userID] = append(userSecrets[:i], userSecrets[i+1:]...)
				return nil
			}
		}
	}
	return NewErrSecretNotFound(secretID)
}
