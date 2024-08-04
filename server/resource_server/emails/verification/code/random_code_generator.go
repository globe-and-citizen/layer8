package code

import (
	"globe-and-citizen/layer8/server/resource_server/models"
	"math/rand"
	"strconv"
	"strings"
)

type RandomCodeGenerator struct {
	verificationCodeSize int
}

func NewRandomCodeGenerator(verificationCodeSize int) *RandomCodeGenerator {
	g := new(RandomCodeGenerator)
	g.verificationCodeSize = verificationCodeSize
	return g
}

func (g *RandomCodeGenerator) GenerateCode(user *models.User, emailAddress string) (string, error) {
	verificationCode := make([]string, g.verificationCodeSize)
	for i := 0; i < g.verificationCodeSize; i++ {
		verificationCode[i] = strconv.Itoa(rand.Intn(10))
	}
	return strings.Join(verificationCode, ""), nil
}
