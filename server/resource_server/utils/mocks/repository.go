package mocks

import (
	"fmt"
	serverModels "globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/resource_server/dto"

	"globe-and-citizen/layer8/server/resource_server/models"
	"time"
)

type MockRepository struct {
	RegisterUserMock                 func(req dto.RegisterUserDTO, hashedPassword string, salt string) error
	RegisterClientMock               func(client models.Client) error
	FindUserMock                     func(userId uint) (models.User, error)
	SaveEmailVerificationDataMock    func(data models.EmailVerificationData) error
	GetEmailVerificationDataMock     func(userId uint) (models.EmailVerificationData, error)
	SaveProofOfEmailVerificationMock func(userID uint, verificationCode string, proof []byte, zkKeyPairId uint) error
	GetUserForUsernameMock           func(username string) (models.User, error)
	UpdateUserPasswordMock           func(username string, password string) error
	RegisterUserPrecheckMock         func(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error
	RegisterUserv2Mock               func(req dto.RegisterUserDTOv2) error
	UpdateUserPasswordV2Mock         func(username string, storedKey string, serverKey string) error
	DeleteEmailVerificationDataMock  func(userId uint) error
	SetUserEmailVerifiedMock         func(userID uint) error
}

func (m *MockRepository) RegisterUser(req dto.RegisterUserDTO, hashedPassword string, salt string) error {
	return m.RegisterUserMock(req, hashedPassword, salt)
}

func (m *MockRepository) RegisterClient(client models.Client) error {
	return m.RegisterClientMock(client)
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

func (m *MockRepository) UpdateUserPassword(username string, password string) error {
	return m.UpdateUserPasswordMock(username, password)
}

func (m *MockRepository) RegisterPrecheckUser(req dto.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	if m.RegisterUserPrecheckMock != nil {
		return m.RegisterUserPrecheckMock(req, salt, iterCount)
	}
	return nil
}

func (m *MockRepository) RegisterUserv2(req dto.RegisterUserDTOv2) error {
	return m.RegisterUserv2Mock(req)
}

func (m *MockRepository) UpdateUserPasswordV2(username string, storedKey string, serverKey string) error {
	return m.UpdateUserPasswordV2Mock(username, storedKey, serverKey)
}

func (m *MockRepository) LoginPreCheckUser(req dto.LoginPrecheckDTO) (string, string, error) {
	return "test_user", "ThisIsARandomSalt123!@#", nil
}

func (m *MockRepository) LoginUser(req dto.LoginUserDTO) (models.User, error) {
	return models.User{
		Username:  "test_user",
		FirstName: "Test",
		LastName:  "User",
		Password:  "34efcb97e704298f3d64159ee858c6c1826755b37523cfac8a79c2130ea7b16f",
		Salt:      "312c4a2c46405ba4f70f7be070f4d4f7cdede09d4b218bf77c01f9706d7505c9",
	}, nil
}

func (m *MockRepository) ProfileUser(userID uint) (models.User, []models.UserMetadata, error) {
	if userID == 1 {
		return models.User{
				Username:  "test_user",
				FirstName: "Test",
				LastName:  "User",
				Password:  "34efcb97e704298f3d64159ee858c6c1826755b37523cfac8a79c2130ea7b16f",
				Salt:      "312c4a2c46405ba4f70f7be070f4d4f7cdede09d4b218bf77c01f9706d7505c9",
			}, []models.UserMetadata{
				{
					UserID: 1,
					Key:    "email_verified",
					Value:  "true",
				},
				{
					UserID: 1,
					Key:    "display_name",
					Value:  "user",
				},
				{
					UserID: 1,
					Key:    "country",
					Value:  "Unknown",
				},
			}, nil
	}
	return models.User{}, []models.UserMetadata{}, fmt.Errorf("User not found")
}

func (m *MockRepository) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return nil
}

func (m *MockRepository) IsBackendURIExists(backendURL string) (bool, error) {
	return true, nil
}

func (m *MockRepository) CheckBackendURI(backendURL string) (bool, error) {
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

func (m *MockRepository) LoginUserPrecheck(username string) (string, error) {
	return "", nil
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

func (m *MockRepository) LoginClient(req dto.LoginClientDTO) (models.Client, error) {
	return models.Client{}, nil
}

func (m *MockRepository) LoginPreCheckClient(req dto.LoginPrecheckDTO) (string, string, error) {
	return "", "", nil
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

func (m *MockRepository) RegisterClientv2(req dto.RegisterClientDTOv2, id string, secret string) error {
	return nil
}

func (m *MockRepository) RegisterPrecheckClient(req dto.RegisterClientPrecheckDTO, salt string, iterCount int) error {
	return nil
}