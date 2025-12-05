package storage

import (
	"fmt"
	"gophkeeper/server/internal/crypto"
	"gophkeeper/server/internal/models"
)

// EncryptedStore wraps a Store and provides transparent encryption/decryption of secret data
type EncryptedStore struct {
	store     Store
	encryptor *crypto.Encryptor
}

// NewEncryptedStore creates a new EncryptedStore that wraps the provided store
// If encryptionKey is empty, encryption is disabled and data passes through unchanged
func NewEncryptedStore(store Store, encryptionKey string) (*EncryptedStore, error) {
	var encryptor *crypto.Encryptor
	var err error

	if encryptionKey != "" {
		encryptor, err = crypto.NewEncryptor(encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create encryptor: %w", err)
		}
	}

	return &EncryptedStore{
		store:     store,
		encryptor: encryptor,
	}, nil
}

// CreateUser delegates to the underlying store (no encryption needed for users)
func (es *EncryptedStore) CreateUser(user models.User) (models.User, error) {
	return es.store.CreateUser(user)
}

// GetUserByLogin delegates to the underlying store (no encryption needed for users)
func (es *EncryptedStore) GetUserByLogin(login string) (models.User, error) {
	return es.store.GetUserByLogin(login)
}

// CreateSecret encrypts the secret data before storing
func (es *EncryptedStore) CreateSecret(secret models.Secret) (models.Secret, error) {
	if es.encryptor != nil && len(secret.Data) > 0 {
		encryptedData, err := es.encryptor.Encrypt(secret.Data)
		if err != nil {
			return models.Secret{}, fmt.Errorf("failed to encrypt secret data: %w", err)
		}
		secret.Data = encryptedData
	}

	return es.store.CreateSecret(secret)
}

// GetSecrets retrieves and decrypts all secrets for a user
func (es *EncryptedStore) GetSecrets(userID int) ([]models.Secret, error) {
	secrets, err := es.store.GetSecrets(userID)
	if err != nil {
		return nil, err
	}

	// Decrypt each secret
	for i := range secrets {
		if es.encryptor != nil && len(secrets[i].Data) > 0 {
			decryptedData, err := es.encryptor.Decrypt(secrets[i].Data)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt secret %d: %w", secrets[i].ID, err)
			}
			secrets[i].Data = decryptedData
		}
	}

	return secrets, nil
}

// GetSecretByID retrieves and decrypts a specific secret
func (es *EncryptedStore) GetSecretByID(userID, secretID int) (models.Secret, error) {
	secret, err := es.store.GetSecretByID(userID, secretID)
	if err != nil {
		return models.Secret{}, err
	}

	// Decrypt the secret data
	if es.encryptor != nil && len(secret.Data) > 0 {
		decryptedData, err := es.encryptor.Decrypt(secret.Data)
		if err != nil {
			return models.Secret{}, fmt.Errorf("failed to decrypt secret: %w", err)
		}
		secret.Data = decryptedData
	}

	return secret, nil
}

// UpdateSecret encrypts the secret data before updating
func (es *EncryptedStore) UpdateSecret(secret models.Secret) (models.Secret, error) {
	if es.encryptor != nil && len(secret.Data) > 0 {
		encryptedData, err := es.encryptor.Encrypt(secret.Data)
		if err != nil {
			return models.Secret{}, fmt.Errorf("failed to encrypt secret data: %w", err)
		}
		secret.Data = encryptedData
	}

	return es.store.UpdateSecret(secret)
}

// DeleteSecret delegates to the underlying store
func (es *EncryptedStore) DeleteSecret(userID, secretID int) error {
	return es.store.DeleteSecret(userID, secretID)
}
