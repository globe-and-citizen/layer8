package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"time"
)

type service struct {
	repository    interfaces.IRepository
	emailVerifier *verification.EmailVerifier

	verificationCodeValidityDuration time.Duration
}

// Newservice creates a new instance of service
func NewService(
	repo interfaces.IRepository,
	emailVerifier *verification.EmailVerifier,
	verificationCodeValidityDuration time.Duration,
) interfaces.IService {
	return &service{
		repository:                       repo,
		emailVerifier:                    emailVerifier,
		verificationCodeValidityDuration: verificationCodeValidityDuration,
	}
}

func (s *service) RegisterUser(req dto.RegisterUserDTO) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if err := validator.New().Struct(req); err != nil {
		return err
	}
	return s.repository.RegisterUser(req)
}

func (s *service) RegisterClient(req dto.RegisterClientDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}

	req.BackendURI = utils.RemoveProtocolFromURL(req.BackendURI)

	return s.repository.RegisterClient(req)
}

func (s *service) GetClientData(clientName string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.GetClientData(clientName)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:          clientData.ID,
		Secret:      clientData.Secret,
		Name:        clientData.Name,
		RedirectURI: clientData.RedirectURI,
		BackendURI:  clientData.BackendURI,
	}
	return clientModel, nil
}

func (s *service) GetClientDataByBackendURL(backendURL string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.GetClientDataByBackendURL(backendURL)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:          clientData.ID,
		Secret:      clientData.Secret,
		Name:        clientData.Name,
		RedirectURI: clientData.RedirectURI,
		BackendURI:  clientData.BackendURI,
	}
	return clientModel, nil
}

func (s *service) LoginPreCheckUser(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	if err := validator.New().Struct(req); err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	username, salt, err := s.repository.LoginPreCheckUser(req)
	if err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	loginPrecheckResp := models.LoginPrecheckResponseOutput{
		Username: username,
		Salt:     salt,
	}
	return loginPrecheckResp, nil
}

func (s *service) LoginPreCheckClient(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	if err := validator.New().Struct(req); err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	username, salt, err := s.repository.LoginPreCheckClient(req)
	if err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	loginPrecheckResp := models.LoginPrecheckResponseOutput{
		Username: username,
		Salt:     salt,
	}
	return loginPrecheckResp, nil
}

func (s *service) LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error) {
	if err := validator.New().Struct(req); err != nil {
		return models.LoginUserResponseOutput{}, err
	}
	user, err := s.repository.LoginUser(req)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}
	tokenResp, err := utils.CompleteLogin(req, user)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}
	return tokenResp, nil
}

func (s *service) LoginClient(req dto.LoginClientDTO) (models.LoginUserResponseOutput, error) {
	if err := validator.New().Struct(req); err != nil {
		return models.LoginUserResponseOutput{}, err
	}
	client, err := s.repository.LoginClient(req)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}

	tokenResp, err := utils.CompleteClientLogin(req, client)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}
	return tokenResp, nil
}

func (s *service) ProfileUser(userID uint) (models.ProfileResponseOutput, error) {
	user, metadata, err := s.repository.ProfileUser(userID)
	if err != nil {
		return models.ProfileResponseOutput{}, err
	}
	profileResp := models.ProfileResponseOutput{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	for _, data := range metadata {
		switch data.Key {
		case "display_name":
			profileResp.DisplayName = data.Value
		case "country":
			profileResp.Country = data.Value
		case "email_verified":
			profileResp.EmailVerified = data.Value == "true"
		}
	}
	return profileResp, nil
}

func (s *service) ProfileClient(userName string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.ProfileClient(userName)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:          clientData.ID,
		Secret:      clientData.Secret,
		Name:        clientData.Username,
		RedirectURI: clientData.RedirectURI,
		BackendURI:  clientData.BackendURI,
	}
	return clientModel, nil
}

func (s *service) VerifyEmail(userID uint) error {
	user, e := s.repository.FindUser(userID)
	if e != nil {
		return e
	}

	verificationCode := s.emailVerifier.GenerateVerificationCode(&user)

	e = s.emailVerifier.SendVerificationEmail(&user, verificationCode)
	if e != nil {
		return e
	}

	e = s.repository.SaveEmailVerificationData(
		models.EmailVerificationData{
			UserId:           user.ID,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(s.verificationCodeValidityDuration).UTC(),
		},
	)

	return e
}

func (s *service) CheckEmailVerificationCode(userId uint, code string) error {
	verificationData, e := s.repository.GetEmailVerificationData(userId)
	if e != nil {
		return e
	}

	e = s.emailVerifier.VerifyCode(&verificationData, code)
	if e != nil {
		return e
	}

	// generate zk-proof of the verification procedure
	proof := "mock_proof"

	e = s.repository.SaveProofOfEmailVerification(userId, code, proof)

	return e
}

func (s *service) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}
	return s.repository.UpdateDisplayName(userID, req)
}
