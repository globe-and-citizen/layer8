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

	emailAsVariables, err := utils.StringToCircuitVariables(email)
	if err != nil {
		t.Fatalf("Unexpected error while converting email to circuit variables: %e", err)
	}

	saltAsVariables, err := utils.StringToCircuitVariables(salt)
	if err != nil {
		t.Fatalf("Unexpected error while converting salt to circuit variables: %e", err)
	}

	codeAsVariables, err := utils.ConvertCodeToCircuitVariables(verificationCode)
	if err != nil {
		t.Fatalf("Unexpected error while converting code to circuit variables: %e", err)
	}

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

	emailAsVariables, err := utils.StringToCircuitVariables(email)
	if err != nil {
		t.Fatalf("Unexpected error while converting email to circuit variables: %e", err)
	}

	saltAsVariables, err := utils.StringToCircuitVariables(salt)
	if err != nil {
		t.Fatalf("Unexpected error while converting salt to circuit variables: %e", err)
	}

	codeAsVariables, err := utils.ConvertCodeToCircuitVariables(verificationCode)
	if err != nil {
		t.Fatalf("Unexpected error while converting code to circuit variables: %e", err)
	}

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
