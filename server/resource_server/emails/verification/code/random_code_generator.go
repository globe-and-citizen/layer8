package code

import (
	"globe-and-citizen/layer8/server/resource_server/constants"
	"math/rand"
	"strconv"
	"strings"
)

type RandomCodeGenerator struct{}

func (g *RandomCodeGenerator) GenerateCode(emailAddress string) string {
	verificationCode := make([]string, constants.VerificationCodeSize)
	for i := 0; i < constants.VerificationCodeSize; i++ {
		verificationCode[i] = strconv.Itoa(rand.Intn(10))
	}
	return strings.Join(verificationCode, "")
}
