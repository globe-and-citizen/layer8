package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-playground/validator/v10"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/sender"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/tokens"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"os"
	"time"
)

type service struct {
	repository     interfaces.IRepository
	emailVerifier  *verification.EmailVerifier
	proofProcessor zk.IProofProcessor
}

// NewService creates a new instance of service
func NewService(
	repo interfaces.IRepository,
	emailVerifier *verification.EmailVerifier,
	proofProcessor zk.IProofProcessor,
) interfaces.IService {
	return &service{
		repository:     repo,
		emailVerifier:  emailVerifier,
		proofProcessor: proofProcessor,
	}
}

func (s *service) RegisterUser(req dto.RegisterUserDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	hashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	return s.repository.RegisterUser(req, hashedAndSaltedPass, rmSalt)
}

func (s *service) RegisterClient(req dto.RegisterClientDTO) error {
	if err := validator.New().Struct(req); err != nil {
		return err
	}

	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)

	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	HashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	req.BackendURI = utils.RemoveProtocolFromURL(req.BackendURI)

	client := models.Client{
		ID:          clientUUID,
		Secret:      clientSecret,
		Name:        req.Name,
		RedirectURI: req.RedirectURI,
		BackendURI:  req.BackendURI,
		Username:    req.Username,
		Password:    HashedAndSaltedPass,
		Salt:        rmSalt,
	}

	return s.repository.RegisterClient(client)
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
		Name:        clientData.Name,
		RedirectURI: clientData.RedirectURI,
		BackendURI:  clientData.BackendURI,
	}
	return clientModel, nil
}

func (s *service) FindUser(userID uint) (models.User, error) {
	return s.repository.FindUser(userID)
}

func (s *service) FindUserForUsername(username string) (models.User, error) {
	return s.repository.FindUserForUsername(username)
}

func (s *service) VerifyEmail(userID uint, userEmail string) error {
	user, e := s.repository.FindUser(userID)
	if e != nil {
		return e
	}

	verificationCode, err := s.emailVerifier.GenerateVerificationCode(&user, userEmail)
	if err != nil {
		return err
	}

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

func (s *service) GenerateZkProofOfEmailVerification(
	user models.User,
	request dto.CheckEmailVerificationCodeDTO,
) ([]byte, error) {
	return s.proofProcessor.GenerateProof(request.Email, user.Salt, request.Code)
}

func (s *service) SaveProofOfEmailVerification(userId uint, verificationCode string, zkProof []byte) error {
	return s.repository.SaveProofOfEmailVerification(userId, verificationCode, zkProof)
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

func (s *service) VerifyUserEmailProof(user *models.User, emailVerificationCode string) error {
	return s.proofProcessor.VerifyProof(emailVerificationCode, user.Salt, user.EmailProof)
}

func (s *service) GeneratePasswordResetToken() ([]byte, error) {
	return tokens.GeneratePasswordResetToken()
}

func (s *service) SavePasswordResetToken(token []byte, user *models.User) error {
	hashedToken := hash(token)

	err := s.repository.SavePasswordResetToken(
		models.PasswordResetTokenData{
			Username:  user.Username,
			Token:     hashedToken[:],
			ExpiresAt: time.Now().Add(10 * time.Minute),
		},
	)

	return err
}

func (s *service) SendPasswordResetToken(tokenBytes []byte, user *models.User, userEmail string) error {
	passwordResetUrl := fmt.Sprintf(
		"%s/api/v1/verify-password-reset-token?token=%s",
		os.Getenv("PROXY_URL"),
		hex.EncodeToString(tokenBytes),
	)

	emailService := sender.NewMailerSendService(os.Getenv("MAILER_SEND_API_KEY"))
	err := emailService.SendEmail(
		&models.Email{
			SenderAddress: fmt.Sprintf(
				"%s@%s",
				os.Getenv("LAYER8_EMAIL_USERNAME"),
				os.Getenv("LAYER8_EMAIL_DOMAIN"),
			),
			RecipientAddress:     userEmail,
			RecipientDisplayName: user.Username,
			Subject:              "Reset your password on the Layer8 portal",
			Content: fmt.Sprintf(
				"Use the following link to reset your password: %s",
				passwordResetUrl,
			),
		},
	)

	return err
}

func (s *service) GetPasswordResetTokenData(token string) (models.PasswordResetTokenData, error) {
	decodedBytes, err := hex.DecodeString(token)
	if err != nil {
		return models.PasswordResetTokenData{}, fmt.Errorf("failed to decode bytes from token: %e", err)
	}

	hashedToken := hash(decodedBytes)

	tokenData, err := s.repository.GetPasswordResetTokenData(hashedToken[:])
	if err != nil {
		return models.PasswordResetTokenData{}, err
	}

	return tokenData, nil
}

func (s *service) ValidatePasswordResetTokenData(tokenData models.PasswordResetTokenData) error {
	if tokenData.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token is expired. Please run the password reset process again")
	}

	return nil
}

func (s *service) UpdateUserPassword(req dto.UpdatePasswordDTO) error {
	user, err := s.repository.FindUserForUsername(req.Username)
	if err != nil {
		return err
	}

	hashedPassword := utils.SaltAndHashPassword(req.NewPassword, user.Salt)

	return s.repository.UpdateUserPassword(req.Username, hashedPassword)
}

func hash(data []byte) [32]byte {
	return sha256.Sum256(data)
}
