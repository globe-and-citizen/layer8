package repository

import (
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
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

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  HashedAndSaltedPass,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Salt:      rmSalt,
	}

	if err := r.connection.Create(&user).Error; err != nil {
		return err
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

	if err := r.connection.Create(&userMetadata).Error; err != nil {
		r.connection.Delete(&user)
		return err
	}

	return nil
}

func (r *Repository) RegisterClient(req dto.RegisterClientDTO) error {

	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)

	client := models.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
	}

	if err := r.connection.Create(&client).Error; err != nil {
		return err
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

func (r *Repository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	var user models.User
	if err := r.connection.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", "", err
	}
	return user.Username, user.Salt, nil
}

func (r *Repository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	var user models.User
	if err := r.connection.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
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

func (r *Repository) VerifyEmail(userID uint) error {
	return r.connection.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", userID, "email_verified").Update("value", "true").Error
}

func (r *Repository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return r.connection.Model(&models.UserMetadata{}).Where("user_id = ? AND key = ?", userID, "display_name").Update("value", req.DisplayName).Error
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
