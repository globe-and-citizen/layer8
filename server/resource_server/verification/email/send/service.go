package send

import "globe-and-citizen/layer8/server/resource_server/models"

type EmailSenderService interface {
	SendEmail(email *models.Email) error
}
