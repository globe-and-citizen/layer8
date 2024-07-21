package code

import (
	"github.com/stretchr/testify/assert"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"testing"
)

const salt = "ajdjsjsaafktyowqqrtgpowrkdkdkfak"

func TestGenerateCode_EmailWithASCIICharacters(t *testing.T) {
	email := "myemail@gmail.com"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(&models.User{Salt: salt}, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}

func TestGenerateCode_EmailWithChineseCharacters(t *testing.T) {
	email := "用户@例子.广告"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(&models.User{Salt: salt}, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}

func TestGenerateCode_EmailWithGermanCharacters(t *testing.T) {
	email := "Dörte@Sörensen.example.com"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(&models.User{Salt: salt}, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}
