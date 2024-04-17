package utils

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/xdg-go/pbkdf2"
)

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}
