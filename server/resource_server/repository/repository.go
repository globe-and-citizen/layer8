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

func (r *Repository) RegisterUser(req dto.RegisterUserDTO) error {
	var user models.User

	tx := r.connection.Begin()
	if err := tx.Where("username = ?", req.Username).First(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not find user: %e", err)
	}

	err := tx.Model(&user).Updates(map[string]interface{}{
		"public_key": req.PublicKey,
		"stored_key": req.StoredKey,
		"server_key": req.ServerKey,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update user: %e", err)
	}

	userMetadata := models.UserMetadata{
		ID:              user.ID,
		IsEmailVerified: false,
		DisplayName:     "",
		Color:           "",
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

func (r *Repository) RegisterClient(req dto.RegisterClientDTO, clientUUID string, clientSecret string) error {
	tx := r.connection.Begin()

	result := tx.Model(&models.Client{}).
		Where("username = ?", req.Username).
		Updates(map[string]interface{}{
			"name":         req.Name,
			"redirect_uri": req.RedirectURI,
			"backend_uri":  req.BackendURI,
			"id":           clientUUID,
			"secret":       clientSecret,
			"stored_key":   req.StoredKey,
			"server_key":   req.ServerKey,
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("could not update client: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no client found with username: %s", req.Username)
	}

	tx.Commit()
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

func (r *Repository) ProfileUser(userID uint) (models.User, models.UserMetadata, error) {
	var user models.User
	if err := r.connection.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, models.UserMetadata{}, err
	}
	var userMetadata models.UserMetadata
	if err := r.connection.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return models.User{}, models.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *Repository) ProfileClient(username string) (models.Client, error) {
	var client models.Client
	if err := r.connection.Where("username = ?", username).First(&client).Error; err != nil {
		return models.Client{}, err
	}
	return client, nil
}

func (r *Repository) SaveProofOfEmailVerification(
	userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(
		&models.User{},
	).Where(
		"id = ?", userId,
	).Updates(map[string]interface{}{
		"verification_code": verificationCode,
		"email_proof":       emailProof,
		"zk_key_pair_id":    zkKeyPairId,
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
		"id = ?", userId,
	).Update("is_email_verified", true).Error

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

func (r *Repository) UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error {
	return r.connection.Model(&models.UserMetadata{}).
		Where("id = ?", userID).
		Updates(models.UserMetadata{
			DisplayName: req.DisplayName,
			Color:       req.Color,
			Bio:         req.Bio,
		}).Error
}

func (r *Repository) SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error) {
	tx := r.connection.Create(&keyPair)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return keyPair.ID, nil
}

func (r *Repository) GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error) {
	var keyPair models.ZkSnarksKeyPair
	err := r.connection.Model(&models.ZkSnarksKeyPair{}).Last(&keyPair).Error

	if err != nil {
		return models.ZkSnarksKeyPair{}, err
	}

	return keyPair, nil
}

func (r *Repository) GetUserForUsername(username string) (models.User, error) {
	var user models.User

	err := r.connection.Model(&models.User{}).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
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

func (r *Repository) RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	user := models.User{
		Username:       req.Username,
		Salt:           salt,
		IterationCount: iterCount,
		PublicKey:      []byte{},
	}

	if err := r.connection.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create a new user: %v", err)
	}

	return nil
}

func (r *Repository) RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error {
	client := models.Client{
		Username:       req.Username,
		Salt:           salt,
		IterationCount: iterCount,
	}

	if err := r.connection.Create(&client).Error; err != nil {
		return fmt.Errorf("failed to create a new client: %v", err)
	}

	return nil
}

func (r *Repository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return r.connection.Model(&models.User{}).
		Where("username=?", username).
		Updates(map[string]interface{}{"stored_key": storedKey, "server_key": serverKey}).Error
}

func (r *Repository) CreateClientTrafficStatisticsEntry(clientId string, rate int) error {
	return r.connection.Create(&models.ClientTrafficStatistics{
		ClientId:                   clientId,
		RatePerByte:                rate,
		TotalUsageBytes:            0,
		UnpaidAmount:               0,
		LastTrafficUpdateTimestamp: time.Now().UTC(),
	}).Error
}

func (r *Repository) GetClientTrafficStatistics(clientId string) (*models.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var clientStatistics models.ClientTrafficStatistics

	err := r.connection.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		return nil, err
	}

	return &clientStatistics, nil
}

func (r *Repository) AddClientTrafficUsage(clientId string, consumedBytes int, now time.Time) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics models.ClientTrafficStatistics
	err := tx.Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	newTrafficBytes := clientStatistics.TotalUsageBytes + consumedBytes
	newUnpaidAmount := clientStatistics.UnpaidAmount + consumedBytes*clientStatistics.RatePerByte

	err = r.connection.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"total_usage_bytes":             newTrafficBytes,
			"unpaid_amount":                 newUnpaidAmount,
			"last_traffic_update_timestamp": now,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) PayClientTrafficUsage(clientId string, amountPaid int) error {
	tx := r.connection.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	var clientStatistics models.ClientTrafficStatistics
	err := tx.Where("client_id = ?", clientId).
		First(&clientStatistics).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if amountPaid < clientStatistics.UnpaidAmount {
		tx.Rollback()
		return fmt.Errorf("full amount must be paid")
	}

	err = r.connection.Model(&models.ClientTrafficStatistics{}).
		Where("client_id = ?", clientId).
		Updates(map[string]interface{}{
			"unpaid_amount": clientStatistics.UnpaidAmount - amountPaid,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) GetAllClientStatistics() ([]models.ClientTrafficStatistics, error) {
	// TODO: is isolation level higher then the default needed?
	var allClientStatistics []models.ClientTrafficStatistics

	err := r.connection.Find(&allClientStatistics).Error
	if err != nil {
		return nil, err
	}

	return allClientStatistics, nil
}
