package code

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

type MIMCCodeGenerator struct{}

func NewMIMCCodeGenerator() *MIMCCodeGenerator {
	return new(MIMCCodeGenerator)
}

func (g *MIMCCodeGenerator) GenerateCode(user *models.User, input string) (string, error) {
	mimcInstance := mimc.NewMiMC()

	frInput, err := utils.StringToFrElements(input)
	if err != nil {
		return "", err
	}

	frSalt, err := utils.StringToFrElements(user.Salt)
	if err != nil {
		return "", err
	}

	encryptedInput := make([]byte, utils.InputFrRepresentationSize*mimc.BlockSize)

	for i := 0; i < utils.InputFrRepresentationSize; i++ {
		inputElementBytes := frInput[i].Bytes()
		saltElementBytes := frSalt[i].Bytes()
		for j := 0; j < fr.Bytes; j++ {
			encryptedInput[i*fr.Bytes+j] = inputElementBytes[j] ^ saltElementBytes[j]
		}
	}

	_, err = mimcInstance.Write(encryptedInput)
	if err != nil {
		return "", fmt.Errorf("error while writing data into the hmac instance: %e", err)
	}

	hashBytes := mimcInstance.Sum(nil)

	return hex.EncodeToString(hashBytes[len(hashBytes)-utils.VerificationCodeSize/2:]), nil
}
