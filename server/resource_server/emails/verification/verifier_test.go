package verification_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"globe-and-citizen/layer8/server/resource_server/emails/sender"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/code"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils/mocks"
	"testing"
	"time"
)

const adminEmail = "layer8@email.com"
const userId uint = 1
const username = "user"
const userEmail = "user@email.com"
const verificationCode = "123456"
const verificationCodeValidityDuration = 2 * time.Minute

var timestamp = time.Date(2024, time.May, 24, 14, 0, 0, 0, time.UTC)
var timestampPlusTwoSeconds = timestamp.Add(time.Second * 2)
var now = func() time.Time {
	return timestamp
}

var mockSenderService sender.EmailService
var mockCodeGenerator code.Generator

func SetUp() {
	mockSenderService = &mocks.MockEmailSenderService{
		SendEmailFunc: func(email *models.Email) error {
			if email.RecipientAddress != userEmail ||
				email.SenderAddress != adminEmail ||
				email.Content.Username != username ||
				email.Content.Code != verificationCode {
				return fmt.Errorf("")
			}
			return nil
		},
	}
	mockCodeGenerator = &mocks.MockCodeGenerator{
		VerificationCode: verificationCode,
	}
}

func TestGenerateVerificationCode(t *testing.T) {
	SetUp()
	verifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)

	generatedCode, err := verifier.GenerateVerificationCode(
		&models.User{
			ID:       userId,
			Username: username,
		},
		userEmail,
	)

	assert.Nil(t, err)
	assert.Equal(t, generatedCode, verificationCode)
}

func TestSendVerificationEmail(t *testing.T) {
	SetUp()
	verifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)

	e := verifier.SendVerificationEmail(
		&models.User{
			ID:       userId,
			Username: username,
		},
		userEmail,
		verificationCode,
	)

	assert.Nil(t, e)
}

func TestVerifyCode_VerificationCodeIsCorrect(t *testing.T) {
	SetUp()
	verifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)

	e := verifier.VerifyCode(
		&models.EmailVerificationData{
			UserId:           userId,
			VerificationCode: verificationCode,
			ExpiresAt:        timestampPlusTwoSeconds,
		},
		verificationCode,
	)

	assert.Nil(t, e)
}

func TestVerifyCode_VerificationCodeIsIncorrect(t *testing.T) {
	SetUp()
	verifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)

	e := verifier.VerifyCode(
		&models.EmailVerificationData{
			UserId:           userId,
			VerificationCode: verificationCode,
			ExpiresAt:        timestampPlusTwoSeconds,
		},
		"567890",
	)

	assert.NotNil(t, e)
}

func TestVerifyCode_VerificationCodeIsExpired(t *testing.T) {
	SetUp()
	now := func() time.Time {
		return timestampPlusTwoSeconds
	}
	verifier := verification.NewEmailVerifier(
		adminEmail,
		mockSenderService,
		mockCodeGenerator,
		verificationCodeValidityDuration,
		now,
	)

	e := verifier.VerifyCode(
		&models.EmailVerificationData{
			UserId:           userId,
			VerificationCode: verificationCode,
			ExpiresAt:        timestamp,
		},
		verificationCode,
	)

	assert.NotNil(t, e)
}
