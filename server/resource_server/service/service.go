package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/code"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type service struct {
	repository     interfaces.IRepository
	emailVerifier  *verification.EmailVerifier
	proofProcessor zk.IProofProcessor
	codeGenerator  code.Generator
}

// NewService creates a new instance of service
func NewService(
	repo interfaces.IRepository,
	emailVerifier *verification.EmailVerifier,
	proofProcessor zk.IProofProcessor,
	codeGenerator code.Generator,
) interfaces.IService {
	return &service{
		repository:     repo,
		emailVerifier:  emailVerifier,
		proofProcessor: proofProcessor,
		codeGenerator:  codeGenerator,
	}
}

func (s *service) RegisterClient(req dto.RegisterClientDTO) error {
	clientUUID := utils.GenerateUUID()
	clientSecret := utils.GenerateSecret(utils.SecretSize)
	req.BackendURI = utils.RemoveProtocolFromURL(req.BackendURI)

	return s.repository.RegisterClient(req, clientUUID, clientSecret)
}

func (s *service) GetClientData(clientName string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.GetClientData(clientName)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:              clientData.ID,
		Secret:          clientData.Secret,
		Name:            clientData.Name,
		RedirectURI:     clientData.RedirectURI,
		BackendURI:      clientData.BackendURI,
		X509Certificate: clientData.X509Certificate,
	}
	return clientModel, nil
}

func (s *service) GetClientDataByBackendURL(backendURL string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.GetClientDataByBackendURL(backendURL)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:              clientData.ID,
		Secret:          clientData.Secret,
		Name:            clientData.Name,
		RedirectURI:     clientData.RedirectURI,
		BackendURI:      clientData.BackendURI,
		X509Certificate: clientData.X509Certificate,
	}
	return clientModel, nil
}

func (s *service) LoginPrecheckUser(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	sNonce := utils.GenerateRandomSalt(utils.SaltSize)

	user, err := s.repository.GetUserForUsername(req.Username)
	if err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	loginPrecheckResp := models.LoginPrecheckResponseOutput{
		Salt:      user.Salt,
		IterCount: user.IterationCount,
		Nonce:     req.CNonce + sNonce,
	}
	return loginPrecheckResp, nil
}

func (s *service) LoginPrecheckClient(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error) {
	sNonce := utils.GenerateRandomSalt(utils.SaltSize)

	client, err := s.repository.ProfileClient(req.Username)
	if err != nil {
		return models.LoginPrecheckResponseOutput{}, err
	}
	loginPrecheckResp := models.LoginPrecheckResponseOutput{
		Salt:      client.Salt,
		IterCount: client.IterationCount,
		Nonce:     req.CNonce + sNonce,
	}
	return loginPrecheckResp, nil
}

func (s *service) LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error) {
	user, err := s.repository.GetUserForUsername(req.Username)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}

	storedKeyBytes, err := hex.DecodeString(user.StoredKey)
	if err != nil {
		return models.LoginUserResponseOutput{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", req.Username, req.CNonce, user.Salt, user.IterationCount, req.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	clientProofBytes, err := hex.DecodeString(req.ClientProof)
	if err != nil {
		return models.LoginUserResponseOutput{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return models.LoginUserResponseOutput{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	clientKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(clientKeyHash[:])
	if clientKeyHashStr != user.StoredKey {
		return models.LoginUserResponseOutput{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(user.ServerKey)
	if err != nil {
		return models.LoginUserResponseOutput{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	tokenString, err := utils.GenerateToken(user)
	if err != nil {
		return models.LoginUserResponseOutput{}, fmt.Errorf("error generating token: %v", err)
	}

	return models.LoginUserResponseOutput{
		ServerSignature: serverSignatureHex,
		Token:           tokenString,
	}, nil
}

func (s *service) LoginClient(req dto.LoginClientDTO) (models.LoginClientResponseOutput, error) {
	client, err := s.repository.ProfileClient(req.Username)
	if err != nil {
		return models.LoginClientResponseOutput{}, err
	}

	storedKeyBytes, err := hex.DecodeString(client.StoredKey)
	if err != nil {
		return models.LoginClientResponseOutput{}, fmt.Errorf("error decoding stored key: %v", err)
	}

	authMessage := fmt.Sprintf("[n=%s,r=%s,s=%s,i=%d,r=%s]", req.Username, req.CNonce, client.Salt, client.IterationCount, req.Nonce)
	authMessageBytes := []byte(authMessage)

	clientSignatureHMAC := hmac.New(sha256.New, storedKeyBytes)
	clientSignatureHMAC.Write(authMessageBytes)
	clientSignature := clientSignatureHMAC.Sum(nil)

	clientProofBytes, err := hex.DecodeString(req.ClientProof)
	if err != nil {
		return models.LoginClientResponseOutput{}, fmt.Errorf("error decoding client proof: %v", err)
	}

	clientKeyBytes, err := utils.XorBytes(clientSignature, clientProofBytes)
	if err != nil {
		return models.LoginClientResponseOutput{}, fmt.Errorf("error performing XOR operation: %v", err)
	}

	clientKeyHash := sha256.Sum256(clientKeyBytes)

	clientKeyHashStr := hex.EncodeToString(clientKeyHash[:])
	if clientKeyHashStr != client.StoredKey {
		return models.LoginClientResponseOutput{}, fmt.Errorf("server failed to authenticate the user")
	}

	serverKeyBytes, err := hex.DecodeString(client.ServerKey)
	if err != nil {
		return models.LoginClientResponseOutput{}, fmt.Errorf("error decoding server key: %v", err)
	}

	serverSignatureHMAC := hmac.New(sha256.New, serverKeyBytes)
	serverSignatureHMAC.Write(authMessageBytes)
	serverSignatureHex := hex.EncodeToString(serverSignatureHMAC.Sum(nil))

	tokenString, err := utils.CompleteClientLoginv2(client)
	if err != nil {
		return models.LoginClientResponseOutput{}, fmt.Errorf("error generating token: %v", err)
	}

	return models.LoginClientResponseOutput{
		ServerSignature: serverSignatureHex,
		Token:           tokenString,
	}, nil
}

func (s *service) ProfileUser(userID uint) (models.ProfileResponseOutput, error) {
	user, metadata, err := s.repository.ProfileUser(userID)
	if err != nil {
		return models.ProfileResponseOutput{}, err
	}

	return models.ProfileResponseOutput{
		Username:            user.Username,
		DisplayName:         metadata.DisplayName,
		Bio:                 metadata.Bio,
		Color:               metadata.Color,
		EmailVerified:       metadata.IsEmailVerified,
		PhoneNumberVerified: metadata.IsPhoneNumberVerified,
	}, nil
}

func (s *service) ProfileClient(username string) (models.ClientResponseOutput, error) {
	clientData, err := s.repository.ProfileClient(username)
	if err != nil {
		return models.ClientResponseOutput{}, err
	}
	clientModel := models.ClientResponseOutput{
		ID:              clientData.ID,
		Secret:          clientData.Secret,
		Name:            clientData.Name,
		RedirectURI:     clientData.RedirectURI,
		BackendURI:      clientData.BackendURI,
		X509Certificate: clientData.X509Certificate,
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

func (s *service) GenerateZkProof(
	user models.User,
	input string, // email or phone number
	verificationCode string,
) ([]byte, uint, error) {
	return s.proofProcessor.GenerateProof(input, user.Salt, verificationCode)
}

func (s *service) SaveProofOfEmailVerification(
	userId uint, verificationCode string, zkProof []byte, zkKeyPairId uint,
) error {
	return s.repository.SaveProofOfEmailVerification(userId, verificationCode, zkProof, zkKeyPairId)
}

func (s *service) UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error {
	return s.repository.UpdateUserMetadata(userID, req)
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

func (s *service) RegisterUserPrecheck(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error) {
	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)

	err := s.repository.RegisterPrecheckUser(req, rmSalt, iterCount)
	if err != nil {
		return "", err
	}

	return rmSalt, nil
}

func (s *service) RegisterClientPrecheck(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error) {
	rmSalt := utils.GenerateRandomSalt(utils.SaltSize)

	err := s.repository.RegisterPrecheckClient(req, rmSalt, iterCount)
	if err != nil {
		return "", err
	}

	return rmSalt, nil
}

func (s *service) RegisterUser(req dto.RegisterUserDTO) error {
	return s.repository.RegisterUser(req)
}

func (s *service) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return s.repository.UpdateUserPassword(username, storedKey, serverKey)
}

func (s *service) GetClientUnpaidAmount(clientId string) (int, error) {
	stat, err := s.repository.GetClientTrafficStatistics(clientId)
	if err != nil {
		log.Printf("repository: %e\n", err)
		return 0, err
	}

	return stat.UnpaidAmount, nil
}

func (s *service) RefreshTelegramMessages(baseURL string, offset int64) ([]dto.MessageUpdateDTO, error) {
	requestUrl := fmt.Sprintf("%s/getUpdates?timeout=2&limit=50", baseURL)
	if offset > 0 {
		requestUrl += "&offset=" + strconv.FormatInt(offset, 10)
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getUpdates request to the Telegram API")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response dto.GetUpdatesResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode getUpdates response from the Telegram API")
	}

	if !response.Ok {
		return nil, fmt.Errorf("received not-ok response from the getUpdates endpoint of Telegram API")
	}

	return response.Result, nil
}

func (s *service) SendTelegramBotMessage(base string, request dto.SendMessageRequestDTO) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal the send message request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/sendMessage", base),
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create the sendMessage request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response dto.SendMessageResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Ok {
		return fmt.Errorf("received not-ok response from the sendMessage endpoint of Telegram API")
	}

	return nil
}

func (s *service) GeneratePhoneNumberVerificationCode(user *models.User, phoneNumber string) (string, error) {
	return s.codeGenerator.GenerateCode(user, phoneNumber)
}

func (s *service) SavePhoneNumberVerificationData(
	userID uint,
	verificationCode string,
	zkProof []byte,
	zkPhoneNumberPairID uint,
) error {
	verificationCodeValidityDuration, err := time.ParseDuration(os.Getenv("VERIFICATION_CODE_VALIDITY_DURATION"))
	if err != nil {
		log.Fatalf("error parsing verification code validity duration: %e", err)
	}

	return s.repository.SavePhoneNumberVerificationData(models.PhoneNumberVerificationData{
		UserId:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().UTC().Add(verificationCodeValidityDuration),
		ZkProof:          zkProof,
		ZkPairID:         zkPhoneNumberPairID,
	})
}

func (s *service) GetPhoneNumberVerificationData(userID uint) (models.PhoneNumberVerificationData, error) {
	return s.repository.GetPhoneNumberVerificationData(userID)
}

func (s *service) CheckPhoneNumberVerificationCode(
	verificationCode string,
	userVerificationData models.PhoneNumberVerificationData,
) error {
	if verificationCode != userVerificationData.VerificationCode {
		return fmt.Errorf("invalid verification code")
	}

	if userVerificationData.ExpiresAt.Before(time.Now().UTC()) {
		return fmt.Errorf("verification code is expired. Try to verify your email again")
	}

	return nil
}

func (s *service) SaveProofOfPhoneNumberVerification(verificationData models.PhoneNumberVerificationData) error {
	return s.repository.SaveProofOfPhoneNumberVerification(
		verificationData.UserId,
		verificationData.VerificationCode,
		verificationData.ZkProof,
		verificationData.ZkPairID,
	)
}

func (s *service) GenerateTelegramSessionID() ([]byte, error) {
	const entropyBytes = 32
	b := make([]byte, entropyBytes)
	if _, err := rand.Read(b); err != nil {
		return []byte{}, fmt.Errorf("read random: %w", err)
	}
	return b, nil
}

func (s *service) SaveTelegramSessionID(userID uint, sessionID []byte) error {
	// hash sessionID first
	sessionIDHash := sha256.Sum256(sessionID)

	return s.repository.SaveTelegramSessionIDHash(userID, sessionIDHash[:])
}
