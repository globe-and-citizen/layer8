package verification

import (
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/emails/sender"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/code"
	"globe-and-citizen/layer8/server/resource_server/models"
	"time"
)

type EmailVerifier struct {
	adminEmailAddress string

	emailSenderService sender.EmailService
	codeGenerator      code.Generator

	now func() time.Time

	VerificationCodeValidityDuration time.Duration
}

func NewEmailVerifier(
	adminEmailAddress string,
	emailSenderService sender.EmailService,
	codeGenerator code.Generator,
	verificationCodeValidityDuration time.Duration,
	now func() time.Time,
) *EmailVerifier {
	verifier := new(EmailVerifier)

	verifier.adminEmailAddress = adminEmailAddress
	verifier.emailSenderService = emailSenderService
	verifier.codeGenerator = codeGenerator
	verifier.now = now

	verifier.VerificationCodeValidityDuration = verificationCodeValidityDuration

	return verifier
}

func (v *EmailVerifier) GenerateVerificationCode(user *models.User, userEmail string) (string, error) {
	return v.codeGenerator.GenerateCode(user, userEmail)
}

func (v *EmailVerifier) SendVerificationEmail(
	user *models.User, userEmail string, verificationCode string,
) error {
	return v.emailSenderService.SendEmail(
		&models.Email{
			SenderAddress:        v.adminEmailAddress,
			RecipientAddress:     userEmail,
			RecipientDisplayName: user.Username,
			Subject:              "Verify your email at the Layer8 service",
			Content: fmt.Sprintf(
				"Hi, %s!\nYour verification code is: %s",
				user.Username,
				verificationCode,
			),
		},
	)
}

func (v *EmailVerifier) VerifyCode(verificationData *models.EmailVerificationData, code string) error {
	if verificationData.ExpiresAt.Before(v.now()) {
		return fmt.Errorf(
			"the verification code is expired. Please try to run the verification process again",
		)
	}

	if code != verificationData.VerificationCode {
		return fmt.Errorf(
			"invalid verification code, expected %s, got %s",
			verificationData.VerificationCode,
			code,
		)
	}

	return nil
}
