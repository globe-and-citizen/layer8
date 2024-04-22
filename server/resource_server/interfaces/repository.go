package interfaces

import (
	serverModel "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"time"
)

type IRepository interface {
	// Resource Server methods
	RegisterUser(req dto.RegisterUserDTO) error
	LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error)
	LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error)
	LoginUser(req dto.LoginUserDTO) (models.User, error)
	LoginClient(req dto.LoginClientDTO) (models.Client, error)
	ProfileUser(userID uint) (models.User, []models.UserMetadata, error)
	ProfileClient(username string) (models.Client, error)
	VerifyEmail(userID uint) error
	UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error
	RegisterClient(req dto.RegisterClientDTO) error
	GetClientData(clientName string) (models.Client, error)
	GetAllClientData() ([]models.Client, error)
	// Oauth2 methods
	LoginUserPrecheck(username string) (string, error)
	GetUser(username string) (*serverModel.User, error)
	GetUserByID(id int64) (*serverModel.User, error)
	GetUserMetadata(userID int64, key string) (*serverModel.UserMetadata, error)
	SetClient(client *serverModel.Client) error
	GetClient(id string) (*serverModel.Client, error)
	SetTTL(key string, value []byte, ttl time.Duration) error
	GetTTL(key string) ([]byte, error)
}
