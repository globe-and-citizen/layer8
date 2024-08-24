package sender

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"globe-and-citizen/layer8/server/resource_server/models"
	"net/http"
	"time"
)

type EmailService interface {
	SendEmail(email *models.Email) error
}

const Layer8EmailDisplayName = "Layer8 team"
const emailSendTimeout = time.Second * 10

type MailerSendService struct {
	apiKey string
}

func NewMailerSendService(apiKey string) *MailerSendService {
	ms := new(MailerSendService)
	ms.apiKey = apiKey
	return ms
}

func (ms *MailerSendService) SendEmail(email *models.Email) error {
	mailerSendClient := mailersend.NewMailersend(ms.apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, emailSendTimeout)
	defer cancel()

	from := mailersend.From{
		Name:  Layer8EmailDisplayName,
		Email: email.SenderAddress,
	}
	to := mailersend.Recipient{
		Name:  email.RecipientDisplayName,
		Email: email.RecipientAddress,
	}

	message := mailerSendClient.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients([]mailersend.Recipient{to})
	message.SetSubject(email.Subject)
	message.SetText(email.Content)

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

	return nil
}
