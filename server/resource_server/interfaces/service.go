package interfaces

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
)

type IService interface {
	RegisterUser(req dto.RegisterUserDTO) error
	LoginPreCheckUser(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error)
	LoginPreCheckClient(req dto.LoginPrecheckDTO) (models.LoginPrecheckResponseOutput, error)
	LoginUser(req dto.LoginUserDTO) (models.LoginUserResponseOutput, error)
	LoginClient(req dto.LoginClientDTO) (models.LoginUserResponseOutput, error)
	ProfileUser(userID uint) (models.ProfileResponseOutput, error)
	ProfileClient(userID string) (models.ClientResponseOutput, error)
	VerifyEmail(userID uint) error
	UpdateDisplayName(userID uint, req dto.UpdateDisplayNameDTO) error
	RegisterClient(req dto.RegisterClientDTO) error
	GetClientData(clientName string) (models.ClientResponseOutput, error)
}
