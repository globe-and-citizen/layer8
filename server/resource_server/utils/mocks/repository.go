package mocks

import (
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"

	"globe-and-citizen/layer8/server/resource_server/models"
	"time"
)

type MockRepository struct {
	RegisterClientMock               func(client models.Client) error
	FindUserMock                     func(userId uint) (models.User, error)
	SaveEmailVerificationDataMock    func(data models.EmailVerificationData) error
	GetEmailVerificationDataMock     func(userId uint) (models.EmailVerificationData, error)
	SaveProofOfEmailVerificationMock func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error
	GetUserForUsernameMock           func(username string) (models.User, error)
	RegisterUserPrecheckMock         func(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error
	RegisterUserMock                 func(req dto.RegisterUserDTO) error
	UpdateUserPasswordMock           func(username string, storedKey string, serverKey string) error
	DeleteEmailVerificationDataMock  func(userId uint) error
	SetUserEmailVerifiedMock         func(userID uint) error
}

func (m *MockRepository) FindUser(userId uint) (models.User, error) {
	return m.FindUserMock(userId)
}

func (m *MockRepository) SaveEmailVerificationData(data models.EmailVerificationData) error {
	return m.SaveEmailVerificationDataMock(data)
}

func (m *MockRepository) GetEmailVerificationData(userId uint) (models.EmailVerificationData, error) {
	return m.GetEmailVerificationDataMock(userId)
}

func (m *MockRepository) GetUserForUsername(username string) (models.User, error) {
	return m.GetUserForUsernameMock(username)
}

func (m *MockRepository) SaveProofOfEmailVerification(
	userID uint, verificationCode string, proof []byte, zkKeyPairId uint,
) error {
	return m.SaveProofOfEmailVerificationMock(userID, verificationCode, proof, zkKeyPairId)
}

func (m *MockRepository) RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	if m.RegisterUserPrecheckMock != nil {
		return m.RegisterUserPrecheckMock(req, salt, iterCount)
	}
	return nil
}

func (m *MockRepository) RegisterUser(req dto.RegisterUserDTO) error {
	return m.RegisterUserMock(req)
}

func (m *MockRepository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return m.UpdateUserPasswordMock(username, storedKey, serverKey)
}

func (m *MockRepository) ProfileUser(userID uint) (models.User, models.UserMetadata, error) {
	return models.User{}, models.UserMetadata{}, nil
}

func (m *MockRepository) UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error {
	return nil
}

func (m *MockRepository) IsBackendURIExists(backendURL string) (bool, error) {
	return true, nil
}

func (m *MockRepository) GetClientData(clientName string) (models.Client, error) {
	if clientName == "testclient" {
		return models.Client{
			ID:          "1",
			Secret:      "testsecret",
			Name:        "testclient",
			RedirectURI: "https://gcitizen.com/callback",
		}, nil
	}
	return models.Client{}, fmt.Errorf("Client not found")
}

func (m *MockRepository) GetUser(username string) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *MockRepository) GetUserByID(id int64) (*serverModels.User, error) {
	return &serverModels.User{}, nil
}

func (m *MockRepository) GetUserMetadata(userID int64, key string) (*serverModels.UserMetadata, error) {
	return &serverModels.UserMetadata{}, nil
}

func (m *MockRepository) SetClient(client *serverModels.Client) error {
	return nil
}

func (m *MockRepository) GetClient(clientName string) (*serverModels.Client, error) {
	return &serverModels.Client{}, nil
}

func (m *MockRepository) SetTTL(key string, value []byte, time time.Duration) error {
	return nil
}

func (m *MockRepository) GetTTL(key string) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockRepository) ProfileClient(userID string) (models.Client, error) {
	return models.Client{}, nil
}

func (m *MockRepository) GetClientDataByBackendURL(backendURL string) (models.Client, error) {
	return models.Client{}, nil
}

func (m *MockRepository) SaveZkSnarksKeyPair(keyPair models.ZkSnarksKeyPair) (uint, error) {
	return 0, nil
}

func (m *MockRepository) GetLatestZkSnarksKeys() (models.ZkSnarksKeyPair, error) {
	return models.ZkSnarksKeyPair{}, nil
}

func (m *MockRepository) RegisterClient(req dto.RegisterClientDTO, id string, secret string) error {
	return nil
}

func (m *MockRepository) RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error {
	return nil
}

func (m *MockRepository) AddClientTrafficUsage(string, int, time.Time) error {
	return nil
}

func (m *MockRepository) CreateClientTrafficStatisticsEntry(string, int) error {
	return nil
}

func (m *MockRepository) GetAllClientStatistics() ([]models.ClientTrafficStatistics, error) {
	return nil, nil
}

func (m *MockRepository) GetClientTrafficStatistics(string) (*models.ClientTrafficStatistics, error) {
	return &models.ClientTrafficStatistics{}, nil
}

func (m *MockRepository) PayClientTrafficUsage(string, int) error {
	return nil
}
