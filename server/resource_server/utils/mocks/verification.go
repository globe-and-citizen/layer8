package mocks

import (
	"globe-and-citizen/layer8/server/resource_server/models"
)

type MockEmailSenderService struct {
	SendEmailFunc func(email *models.Email) error
}

func (s *MockEmailSenderService) SendEmail(email *models.Email) error {
	return s.SendEmailFunc(email)
}

type MockCodeGenerator struct {
	VerificationCode string
}

func (g *MockCodeGenerator) GenerateCode(user *models.User, emailAddress string) (string, error) {
	return g.VerificationCode, nil
}

type MockProofGenerator struct {
	GenerateProofFunc func(emailAddress string, salt string, verificationCode string) ([]byte, error)
	VerifyProofFunc   func(verificationCode string, salt string, proofBytes []byte) error
}

func (pg *MockProofGenerator) GenerateProof(
	emailAddress string, salt string, verificationCode string,
) ([]byte, error) {
	return pg.GenerateProofFunc(emailAddress, salt, verificationCode)
}

func (pg *MockProofGenerator) VerifyProof(
	verificationCode string, salt string, proofBytes []byte,
) error {
	return pg.VerifyProofFunc(verificationCode, salt, proofBytes)
}
