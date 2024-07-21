package zk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const email = "myemail@gmail.com"
const salt = "ajdjsjsaafktyowqqrtgpowrkdkdkfak"

func TestGenerateProof_VerificationCodeIsInvalid(t *testing.T) {
	verificationCode := "724b20"

	zkProofProcessor := NewProofProcessor()
	_, err := zkProofProcessor.GenerateProof(email, salt, verificationCode)

	assert.NotNil(t, err)
}

func TestVerifyProof_ProofIsInvalid(t *testing.T) {
	proof := make([]byte, 164)
	for i := 0; i < len(proof); i++ {
		proof[i] = byte(i)
	}

	zkProofProcessor := NewProofProcessor()
	err := zkProofProcessor.VerifyProof("123456", salt, proof)

	assert.NotNil(t, err)
}

func TestGenerateProof_Success(t *testing.T) {
	verificationCode := "724b2c"

	zkProofProcessor := NewProofProcessor()
	proof, err := zkProofProcessor.GenerateProof(email, salt, verificationCode)

	assert.Nil(t, err)
	assert.True(t, len(proof) > 0)

	err = zkProofProcessor.VerifyProof(verificationCode, salt, proof)
	assert.Nil(t, err)
}
