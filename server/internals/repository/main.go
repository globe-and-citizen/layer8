package repository

import (
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/models"
	"time"
)

type Repository interface {

	// Get the salt from db using the username
	LoginUserPrecheck(username string) (string, error)

	// Get user from db by username
	GetUser(username string) (*models.User, error)

	// GetUserByID gets a user by ID.
	GetUserByID(id int64) (*models.User, error)

	// GetUserMetadata gets a user metadata by key.
	GetUserMetadata(userID int64, key string) (*models.UserMetadata, error)

	// Set a client for testing purposes
	SetClient(client *models.Client) error

	// Get a client by ID.
	GetClient(id string) (*models.Client, error)

	GetClientByURL(url string) (*models.Client, error)

	// SaveX509Certificate saves the x509 certificate per client
	SaveX509Certificate(clientID string, certificate string) error

	// SetTTL sets the value for the given key with a short TTL.
	SetTTL(key string, value []byte, ttl time.Duration) error

	// GetTTL gets the value for the given key which has a short TTL.
	GetTTL(key string) ([]byte, error)
}

func InitDB() Repository {
	// Register the postgres repository
	return NewOauthRepository(config.DB)
}
