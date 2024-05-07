package email

import (
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/verification/email/code"
	"globe-and-citizen/layer8/server/resource_server/verification/email/send"
	"strconv"
	"time"
)

const verificationCodeValidityDuration = time.Minute * 2

type Verifier struct {
	adminEmailAddress string

	emailSenderService send.EmailSenderService
	codeGenerator      code.Generator
}

func (v *Verifier) InitVerification(user *models.User) (models.EmailVerificationData, error) {
	verificationCode := v.codeGenerator.GenerateCode(user.Email)
	expiresAt := time.Now().Add(verificationCodeValidityDuration)

	e := v.emailSenderService.SendEmail(
		&models.Email{
			From:    v.adminEmailAddress,
			To:      user.Email,
			Subject: "Verify your email at the Layer8 service",
			Content: models.VerificationEmailContent{
				Username: user.Username,
				Code:     verificationCode,
			},
		},
	)

	if e != nil {
		return models.EmailVerificationData{}, e
	}

	verificationData :=
		models.EmailVerificationData{
			UserId:           strconv.Itoa(int(user.ID)),
			VerificationCode: verificationCode,
			ExpiresAt:        expiresAt,
		}
	return verificationData, nil
}

func (v *Verifier) VerifyCode(code string, verificationData models.EmailVerificationData) error {
	if verificationData.ExpiresAt.Before(time.Now()) {
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
