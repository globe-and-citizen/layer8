package interfaces

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
)

type IService interface {
	RegisterUser(req dto.RegisterUserDTO) error
	LoginPrecheckUser(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error)
	LoginPrecheckClient(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error)
	LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error)
	LoginClient(req dto.LoginClientDTO) (models.LoginClientResponseOutput, error)
	ProfileUser(userID uint) (models.ProfileResponseOutput, error)
	ProfileClient(userID string) (models.ClientResponseOutput, error)
	FindUser(userID uint) (models.User, error)
	VerifyEmail(userID uint, userEmail string) error
	CheckEmailVerificationCode(userID uint, code string) error
	GenerateZkProofOfEmailVerification(
		user models.User, request dto.CheckEmailVerificationCodeDTO,
	) ([]byte, uint, error)
	SaveProofOfEmailVerification(userID uint, verificationCode string, zkProof []byte, zkKeyPairId uint) error
	UpdateUserMetadata(userID uint, req dto.UpdateUserMetadataDTO) error
	RegisterClient(req dto.RegisterClientDTO) error
	GetClientData(clientName string) (models.ClientResponseOutput, error)
	GetClientDataByBackendURL(backendURL string) (models.ClientResponseOutput, error)
	CheckBackendURI(backendURL string) (bool, error)
	GetUserForUsername(username string) (models.User, error)
	ValidateSignature(message string, signature []byte, publicKey []byte) error
	UpdateUserPassword(username string, storedKey string, serverKey string) error
	RegisterUserPrecheck(req dto.RegisterUserPrecheckDTO, iterCount int) (string, error)
	RegisterClientPrecheck(req dto.RegisterClientPrecheckDTO, iterCount int) (string, error)
	GetClientUnpaidAmount(clientId string) (int, error)
}
