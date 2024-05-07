package send

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"globe-and-citizen/layer8/server/resource_server/models"
	"net/http"
	"time"
)

const emailSendTimeout = time.Second * 10

type MailerSendService struct {
	apiKey     string
	templateId string
}

func NewMailerSendService(apiKey string, templateId string) *MailerSendService {
	ms := new(MailerSendService)
	ms.apiKey = apiKey
	ms.templateId = templateId
	return ms
}

func (ms *MailerSendService) SendEmail(email *models.Email) error {
	mailerSendClient := mailersend.NewMailersend(ms.apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, emailSendTimeout)
	defer cancel()

	from := mailersend.From{
		Name:  "Layer8 team",
		Email: email.From,
	}
	to := mailersend.Recipient{
		Name:  email.Content.Username,
		Email: email.To,
	}

	personalization := []mailersend.Personalization{
		{
			Email: email.To,
			Data: map[string]interface{}{
				"code": email.Content.Code,
				"user": email.Content.Username,
			},
		},
	}

	message := mailerSendClient.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients([]mailersend.Recipient{to})
	message.SetSubject(email.Subject)
	message.SetTemplateID(ms.templateId)
	message.SetPersonalization(personalization)

	response, e := mailerSendClient.Email.Send(ctx, message)
	if e != nil {
		return fmt.Errorf("error while sending a verification email via MailerSend: %e", e)
	}
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"failed to send a verification email, status code %d",
			response.StatusCode,
		)
	}

	return e
}
