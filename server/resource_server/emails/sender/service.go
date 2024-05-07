package sender

import "globe-and-citizen/layer8/server/resource_server/models"

type Service interface {
	SendEmail(email *models.Email) error
}
