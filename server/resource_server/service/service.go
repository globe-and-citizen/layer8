package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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
	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)
	hashedAndSaltedPass := utils.SaltAndHashPassword(req.Password, rmSalt)

	return s.repository.RegisterUser(req, hashedAndSaltedPass, rmSalt)
}

func (s *service) RegisterClient(req dto.RegisterClientDTO) error {
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
) ([]byte, uint, error) {
	return s.proofProcessor.GenerateProof(request.Email, user.Salt, request.Code)
}

func (s *service) SaveProofOfEmailVerification(
	userId uint, verificationCode string, zkProof []byte, zkKeyPairId uint,
) error {
	return s.repository.SaveProofOfEmailVerification(userId, verificationCode, zkProof, zkKeyPairId)
}

func (s *service) UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error {
	return s.repository.UpdateDisplayName(userID, req)
}

func (s *service) CheckBackendURI(backendURL string) (bool, error) {
	response, err := s.repository.IsBackendURIExists(backendURL)
	if err != nil {
		return false, err
	}
	return response, nil
}

func (s *service) GetUserForUsername(username string) (models.User, error) {
	return s.repository.GetUserForUsername(username)
}

func (s *service) ValidateSignature(message string, signature []byte, publicKey []byte) error {
	msgHash := crypto.Keccak256([]byte(message))
	verified := crypto.VerifySignature(publicKey, msgHash, signature)

	if !verified {
		return fmt.Errorf("failed to verify the ecdsa signature")
	}

	return nil
}

func (s *service) UpdateUserPassword(username string, newPassword string, salt string) error {
	hashedPassword := utils.SaltAndHashPassword(newPassword, salt)
	return s.repository.UpdateUserPassword(username, hashedPassword)
}

func (s *service) RegisterUserPrecheck(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)

	err := s.repository.RegisterPrecheckUser(req, rmSalt, iterCount)
	if err != nil {
		return "", err
	}

	return rmSalt, nil
}

func (s *service) RegisterUserv2(req dto.RegisterUserDTOv2) error {
	return s.repository.RegisterUserv2(req)
}

func (s *service) LoginPreCheckUserv2(req dto.LoginPrecheckDTOv2) (models.LoginPrecheckResponseOutputv2, error) {
	sNonce := utils.GenerateRandomSalt(32)

	user, err := s.repository.GetUserForUsername(req.Username)
	if err != nil {
		return models.LoginPrecheckResponseOutputv2{}, err
	}
	loginPrecheckResp := models.LoginPrecheckResponseOutputv2{
		Salt:      user.Salt,
		IterCount: user.IterationCount,
		Nonce:     req.CNonce + sNonce,
	}
	return loginPrecheckResp, nil
}

func (s *service) LoginUserv2(req dto.LoginUserDTOv2) (models.LoginUserResponseOutputv2, error) {
	user, err := s.repository.GetUserForUsername(req.Username)
	if err != nil {
		return models.LoginUserResponseOutputv2{}, err
	}

	storedKeyBytes, err := hex.DecodeString(user.StoredKey)
	if err != nil {
		return models.LoginUserResponseOutputv2{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", req.Username, req.CNonce, user.Salt, user.IterationCount, req.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	fmt.Println("Client Signature Bytes: ", clientSignature)

	clientProofBytes, err := utils.HexStringToBytes(req.ClientProof)
	if err != nil {
		return models.LoginUserResponseOutputv2{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	fmt.Println("Client Proof: ", clientProofBytes)

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return models.LoginUserResponseOutputv2{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	fmt.Println("Client Key Bytes: ", clientKeyBytes)

	cleintKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(cleintKeyHash[:])

	fmt.Println("Client Key SHA256 String: ", clientKeyHashStr)
	fmt.Println("Stored Key: ", user.StoredKey)

	if clientKeyHashStr != user.StoredKey {
		return models.LoginUserResponseOutputv2{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(user.ServerKey)
	if err != nil {
		return models.LoginUserResponseOutputv2{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	// For testing only
	fmt.Println("Server Signature:", serverSignatureHex)

	return models.LoginUserResponseOutputv2{
		ServerSignature: serverSignatureHex,
	}, nil
}
