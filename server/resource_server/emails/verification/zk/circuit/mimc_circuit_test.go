package circuit

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"testing"
)

const email = "myemail@gmail.com"
const salt = "ajdjsjsaafktyowqqrtgpowrkdkdkfak"

func TestMimcCircuit_ProverSucceeded(t *testing.T) {
	verificationCode := "724b2c"

	emailAsVariables, _ := utils.StringToCircuitVariables(email)
	saltAsVariables, _ := utils.StringToCircuitVariables(salt)
	codeAsVariables, _ := utils.ConvertCodeToCircuitVariables(verificationCode)

	circuit := NewMimcCircuit()

	test.NewAssert(t).ProverSucceeded(
		circuit,
		&MimcCircuit{
			EmailAsVariables: emailAsVariables,
			SaltAsVariables:  saltAsVariables,
			VerificationCode: codeAsVariables,
		},
		test.WithCurves(ecc.BN254),
	)
}

func TestMimcCircuit_ProverFailed(t *testing.T) {
	verificationCode := "724b2b"

	emailAsVariables, _ := utils.StringToCircuitVariables(email)
	saltAsVariables, _ := utils.StringToCircuitVariables(salt)
	codeAsVariables, _ := utils.ConvertCodeToCircuitVariables(verificationCode)

	circuit := NewMimcCircuit()

	test.NewAssert(t).ProverFailed(
		circuit,
		&MimcCircuit{
			EmailAsVariables: emailAsVariables,
			SaltAsVariables:  saltAsVariables,
			VerificationCode: codeAsVariables,
		},
		test.WithCurves(ecc.BN254),
	)
}
