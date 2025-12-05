package storage

import "gophkeeper/server/internal/models"

type Store interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByLogin(login string) (models.User, error)

	CreateSecret(secret models.Secret) (models.Secret, error)
	GetSecrets(userID int) ([]models.Secret, error)
	GetSecretByID(userID, secretID int) (models.Secret, error)
	UpdateSecret(secret models.Secret) (models.Secret, error)
	DeleteSecret(userID, secretID int) error
}
