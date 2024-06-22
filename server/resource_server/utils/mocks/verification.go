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

func (g *MockCodeGenerator) GenerateCode(emailAddress string) string {
	return g.VerificationCode
}
