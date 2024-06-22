package service

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/server/blockchain"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repository        interfaces.IRepository
	emailVerifier     *verification.EmailVerifier
	payAsYouGoWrapper blockchain.PayAsYouGoWrapper
}

// Newservice creates a new instance of service
func NewService(
	repo interfaces.IRepository,
	emailVerifier *verification.EmailVerifier,
	payAsYouGoWrapper blockchain.PayAsYouGoWrapper,
) interfaces.IService {
	return &service{
		repository:        repo,
		emailVerifier:     emailVerifier,
		payAsYouGoWrapper: payAsYouGoWrapper,
	}
}

func (s *service) RegisterUser(req dto.RegisterUserDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}
	return s.repository.RegisterUser(req)
}

func (s *service) RegisterClient(req dto.RegisterClientDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}

	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	client := &models.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
		BackendURI:  utils.RemoveProtocolFromURL(req.BackendURI),
		Username:    req.Username,
		Password:    HashedAndSaltedPass,
		Salt:        rmSalt,
	}

	if err := s.repository.RegisterClient(client); err != nil {
		return err
	}

	go func() {
		billRate, err := strconv.Atoi(os.Getenv("BLOCKCHAIN_BILL_RATE"))
		if err != nil {
			fmt.Println("Error converting bill rate to int: ", err)
			return
		}

		contractID, err := s.payAsYouGoWrapper.CreateContract(context.Background(), uint8(billRate), clientUUID)
		if err != nil {
			fmt.Println("Error creating contract: ", err)
			return
		}

		err = s.repository.UpdateClientBlockchainContractID(clientUUID, contractID)
		if err != nil {
			fmt.Println("Error updating client contract ID: ", err)
			return
		}

	}()

	return nil
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

func (s *service) VerifyEmail(userID uint, userEmail string) error {
	user, e := s.repository.FindUser(userID)
	if e != nil {
		return e
	}

	verificationCode := s.emailVerifier.GenerateVerificationCode(&user, userEmail)

	e = s.emailVerifier.SendVerificationEmail(&user, userEmail, verificationCode)
	if e != nil {
		return e
	}

	e = s.repository.SaveEmailVerificationData(
		models.EmailVerificationData{
			UserId:           user.ID,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(s.emailVerifier.VerificationCodeValidityDuration).UTC(),
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

	return e
}

func (s *service) GenerateZkProofOfEmailVerification(userID uint) (string, error) {
	return "mock_proof", nil
}

func (s *service) SaveProofOfEmailVerification(userId uint, verificationCode string, zkProof string) error {
	e := s.repository.SaveProofOfEmailVerification(userId, verificationCode, zkProof)

	return e
}

func (s *service) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}
	return s.repository.UpdateDisplayName(userID, req)
}

func (s *service) CheckBackendURI(backendURL string) (bool, error) {
	response, err := s.repository.IsBackendURIExists(backendURL)
	if err != nil {
		return false, err
	}
	return response, nil
}
