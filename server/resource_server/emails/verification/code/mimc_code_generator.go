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

func (g *MIMCCodeGenerator) GenerateCode(user *models.User, emailAddress string) (string, error) {
	mimcInstance := mimc.NewMiMC()

	frEmail, err := utils.StringToFrElements(emailAddress)
	if err != nil {
		return "", err
	}

	frSalt, err := utils.StringToFrElements(user.Salt)
	if err != nil {
		return "", err
	}

	encryptedEmail := make([]byte, utils.EmailFrRepresentationSize*mimc.BlockSize)

	for i := 0; i < utils.EmailFrRepresentationSize; i++ {
		emailElementBytes := frEmail[i].Bytes()
		saltElementBytes := frSalt[i].Bytes()
		for j := 0; j < fr.Bytes; j++ {
			encryptedEmail[i*fr.Bytes+j] = emailElementBytes[j] ^ saltElementBytes[j]
		}
	}

	_, err = mimcInstance.Write(encryptedEmail)
	if err != nil {
		return "", fmt.Errorf("error while writing data into the hmac instance: %e", err)
	}

	hashBytes := mimcInstance.Sum(nil)

	return hex.EncodeToString(hashBytes[len(hashBytes)-utils.VerificationCodeSize/2:]), nil
}
