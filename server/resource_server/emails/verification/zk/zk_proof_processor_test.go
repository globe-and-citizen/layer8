package zk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const email = "myemail@gmail.com"
const salt = "ajdjsjsaafktyowqqrtgpowrkdkdkfak"
const zkKeyPairId uint = 2

func TestGenerateProof_VerificationCodeIsInvalid(t *testing.T) {
	verificationCode := "724b20"

	cs, provingKey, verifyingKey := RunZkSnarksSetup()
	zkProofProcessor := NewProofProcessor(cs, zkKeyPairId, provingKey, verifyingKey)
	_, _, err := zkProofProcessor.GenerateProof(email, salt, verificationCode)

	assert.NotNil(t, err)
}

func TestGenerateProof_Success(t *testing.T) {
	verificationCode := "724b2c"

	cs, provingKey, verifyingKey := RunZkSnarksSetup()
	zkProofProcessor := NewProofProcessor(cs, zkKeyPairId, provingKey, verifyingKey)
	proof, actualZkKeyPairId, err := zkProofProcessor.GenerateProof(email, salt, verificationCode)

	assert.Nil(t, err)
	assert.Equal(t, zkKeyPairId, actualZkKeyPairId)
	assert.True(t, len(proof) > 0)

	err = zkProofProcessor.VerifyProof(verificationCode, salt, proof)
	assert.Nil(t, err)
}

func TestVerifyProof_ProofIsInvalid(t *testing.T) {
	proof := make([]byte, 164)
	for i := 0; i < len(proof); i++ {
		proof[i] = byte(i)
	}

	cs, provingKey, verifyingKey := RunZkSnarksSetup()
	zkProofProcessor := NewProofProcessor(cs, zkKeyPairId, provingKey, verifyingKey)
	err := zkProofProcessor.VerifyProof("123456", salt, proof)

	assert.NotNil(t, err)
}
