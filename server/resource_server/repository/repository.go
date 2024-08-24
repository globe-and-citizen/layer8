package repository

import (
	"database/sql"
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	connection *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.IRepository {
	return &Repository{
		connection: db,
	}
}

func (r *Repository) RegisterUser(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
	user := models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Salt:      salt,
	}

	tx := r.connection.Begin()

	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create user: %e", err)
	}

	userMetadata := []models.UserMetadata{
		{
			UserID: user.ID,
			Key:    "email_verified",
			Value:  "false",
		},
		{
			UserID: user.ID,
			Key:    "country",
			Value:  req.Country,
		},
		{
			UserID: user.ID,
			Key:    "display_name",
			Value:  req.DisplayName,
		},
	}

	err = tx.Create(&userMetadata).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create user metadata entry: %e", err)
	}

	tx.Commit()

	return nil
}

func (r *Repository) FindUser(userId uint) (models.User, error) {
	var user models.User
	e := r.connection.Where("id = ?", userId).First(&user).Error

	if e != nil {
		return models.User{}, e
	}

	return user, e
}

func (r *Repository) FindUserForUsername(username string) (models.User, error) {
	var user models.User
	err := r.connection.Where("username = ?", username).First(&user).Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) RegisterClient(client models.Client) error {
	if err := r.connection.Create(&client).Error; err != nil {
		return fmt.Errorf("failed to create a new client record: %e", err)
	}

	return nil
}

func (r *Repository) GetClientData(clientName string) (models.Client, error) {
	var client models.Client
	if err := r.connection.Where("name = ?", clientName).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) GetClientDataByBackendURL(backendURL string) (models.Client, error) {
	var client models.Client
	if err := r.connection.Where("backend_uri = ?", backendURL).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	var user models.User
	if err := r.connection.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Salt, nil
}

func (r *Repository) LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error) {
	var client models.Client
	// RAVI
	// if err := config.DB.Where("username = ?", req.Username).First(&client).Error; err != nil {
	// 	return "", "", err
	// }
	if err := r.connection.Where("username = ?", req.Username).First(&client).Error; err != nil {
		return "", "", err
	}
	return client.Username, client.Salt, nil
}

func (r *Repository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	var user models.User
	if err := r.connection.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *Repository) LoginClient(req dto.LoginClientDTO) (models.Client, error) {
	var client models.Client
	// RAVI HERE
	// if err := config.DB.Where("username = ?", req.Username).First(&client).Error; err != nil {
	// 	return models.Client{}, err
	// }
	if err := r.connection.Where("username = ?", req.Username).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) ProfileUser(userID uint) (models.User, []models.UserMetadata, error) {
	var user models.User
	if err := r.connection.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, []models.UserMetadata{}, err
	}
	var userMetadata []models.UserMetadata
	if err := r.connection.Where("user_id = ?", userID).Find(&userMetadata).Error; err != nil {
		return models.User{}, []models.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *Repository) ProfileClient(userID string) (models.Client, error) {
	var client models.Client
	if err := r.connection.Where("username = ?", userID).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) SaveProofOfEmailVerification(
	userId uint, verificationCode string, emailProof []byte,
) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(
		&models.User{},
	).Where(
		"id = ?", userId,
	).Updates(map[string]interface{}{
		"verification_code": verificationCode,
		"email_proof":       emailProof,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(
		"user_id = ?", userId,
	).Delete(
		&models.EmailVerificationData{},
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(
		&models.UserMetadata{},
	).Where(
		"user_id = ? AND key = ?",
		userId,
		"email_verified",
	).Update("value", "true").Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) SaveEmailVerificationData(data models.EmailVerificationData) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(
		models.EmailVerificationData{UserId: data.UserId},
	).Assign(data).FirstOrCreate(&models.EmailVerificationData{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) GetEmailVerificationData(userId uint) (models.EmailVerificationData, error) {
	var data models.EmailVerificationData
	e := r.connection.Where("user_id = ?", userId).First(&data).Error
	if e != nil {
		return models.EmailVerificationData{}, e
	}

	return data, nil
}

func (r *Repository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return r.connection.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", userID, "display_name").Update("value", req.DisplayName).Error
}

func (r *Repository) SavePasswordResetToken(tokenEntry models.PasswordResetTokenData) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	tx.Where(
		models.PasswordResetTokenData{Username: tokenEntry.Username},
	).Assign(tokenEntry).FirstOrCreate(&models.PasswordResetTokenData{})

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (r *Repository) GetPasswordResetTokenData(token []byte) (models.PasswordResetTokenData, error) {
	var tokenRow models.PasswordResetTokenData

	err := r.connection.Where("token = ?", token).First(&tokenRow).Error
	if err != nil {
		return models.PasswordResetTokenData{}, err
	}

	return tokenRow, nil
}

func (r *Repository) UpdateUserPassword(username string, password string) error {
	return r.connection.Model(&models.User{}).
		Where("username = ?", username).
		Update("password", password).
		Error
}

func (r *Repository) SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) error {
	return r.connection.Create(&keyPair).Error
}

func (r *Repository) GetZkSnarksKeys() (models.ZkSnarksKeyPair, error) {
	var keyPair models.ZkSnarksKeyPair
	err := r.connection.Model(&models.ZkSnarksKeyPair{}).Last(&keyPair).Error

	if err != nil {
		return models.ZkSnarksKeyPair{}, err
	}

	return keyPair, nil
}

func (r *Repository) LoginUserPrecheck(username string) (string, error) {
	return "", nil
}

func (r *Repository) GetUser(username string) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (r *Repository) GetUserByID(id int64) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (r *Repository) GetUserMetadata(userID int64, key string) (*serverModels.UserMetadata, error) {
	return &serverModels.UserMetadata{}, nil
}

func (r *Repository) SetClient(client *serverModels.Client) error {
	return nil
}

func (r *Repository) GetClient(id string) (*serverModels.Client, error) {
	return &serverModels.Client{}, nil
}

func (r *Repository) SetTTL(key string, value []byte, time time.Duration) error {
	return nil
}

func (r *Repository) GetTTL(key string) ([]byte, error) {
	return []byte{}, nil
}

func (r *Repository) IsBackendURIExists(backendURL string) (bool, error) {
	var count int64
	if err := r.connection.Model(&models.Client{}).Where("backend_uri = ?", backendURL).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
