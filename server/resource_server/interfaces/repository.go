package interfaces

import (
	serverModel "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"time"
)

type IRepository interface {
	// Resource Server methods
	FindUser(userId uint) (models.User, error)
	ProfileUser(userID uint) (models.User, models.UserMetadata, error)
	ProfileClient(username string) (models.Client, error)
	SaveProofOfEmailVerification(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error
	SaveEmailVerificationData(data models.EmailVerificationData) error
	GetEmailVerificationData(userId uint) (models.EmailVerificationData, error)
	UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error
	RegisterUser(req dto.RegisterUserDTO) error
	RegisterClient(req dto.RegisterClientDTO, clientUUID string, clientSecret string) error
	GetClientData(clientName string) (models.Client, error)
	GetClientDataByBackendURL(backendURL string) (models.Client, error)
	IsBackendURIExists(backendURL string) (bool, error)
	SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error)
	GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error)
	GetUserForUsername(username string) (models.User, error)
	UpdateUserPassword(username string, storedKey string, serverKey string) error
	CreateClientTrafficStatisticsEntry(clientId string, rate int) error
	AddClientTrafficUsage(clientId string, consumedBytes int, now time.Time) error
	GetClientTrafficStatistics(clientId string) (*models.ClientTrafficStatistics, error)
	PayClientTrafficUsage(clientId string, amountPaid int) error
	GetAllClientStatistics() ([]models.ClientTrafficStatistics, error)
	RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error
	RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error

	// Oauth2 methods
	GetUser(username string) (*serverModel.User, error)
	GetUserByID(id int64) (*serverModel.User, error)
	GetUserMetadata(userID int64, key string) (*serverModel.UserMetadata, error)
	SetClient(client *serverModel.Client) error
	GetClient(id string) (*serverModel.Client, error)
	SetTTL(key string, value []byte, ttl time.Duration) error
	GetTTL(key string) ([]byte, error)
}
