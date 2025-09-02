package repository

import (
	"fmt"
	"globe-and-citizen/layer8/server/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db      *gorm.DB
	storage map[string][]byte
}

func NewOauthRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		storage: make(map[string][]byte),
	}
}

func (r *PostgresRepository) LoginUserPrecheck(username string) (string, error) {
	var user models.User
	fmt.Println("RAVI 2: ", user)
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Salt, nil
}

func (r *PostgresRepository) GetUser(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (r *PostgresRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (r *PostgresRepository) GetUserMetadata(userID int64) (*models.UserMetadata, error) {
	var userMetadata models.UserMetadata
	err := r.db.Where("id = ?", userID).First(&userMetadata).Error
	if err != nil {
		return &models.UserMetadata{}, err
	}
	return &userMetadata, nil
}

func (r *PostgresRepository) SetClient(client *models.Client) error {
	// Check if client already exists
	var existingClient models.Client
	err := r.db.Where("id = ?", client.ID).First(&existingClient).Error
	if err == nil {
		return nil
	}
	err = r.db.Create(&client).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetClient(id string) (*models.Client, error) {
	if func() bool {
		var prefix string = "client:"
		return len(id) >= len(prefix) && id[0:len(prefix)] == prefix
	}() {
		id = strings.TrimPrefix(id, "client:")
	}
	var client models.Client
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return &models.Client{}, err
	}
	return &client, nil
}

func (r *PostgresRepository) GetClientByURL(url string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("backend_uri = ?", url).First(&client).Error
	if err != nil {
		return &models.Client{}, err
	}
	return &client, nil
}

func (r *PostgresRepository) SaveX509Certificate(clientID string, certificate string) error {
	return r.db.Model(&models.Client{}).
		Where("id = ?", clientID).
		Update("x509_certificate_bytes", certificate).
		Error
}

// SetTTL sets the key to hold the value for a limited time
func (r *PostgresRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
	r.storage[key] = value
	go func() {
		time.Sleep(ttl)
		delete(r.storage, key)
	}()
	return nil
}

// GetTTL gets the value for the given key which has a short TTL.
func (r *PostgresRepository) GetTTL(key string) ([]byte, error) {
	return r.storage[key], nil
}
