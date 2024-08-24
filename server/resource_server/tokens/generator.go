package tokens

import (
	"crypto/rand"
	"fmt"
)

const PasswordResetTokenSize = 64

func GeneratePasswordResetToken() ([]byte, error) {
	token := make([]byte, PasswordResetTokenSize)

	_, err := rand.Read(token)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password reset token: %e", err)
	}

	return token, nil
}
