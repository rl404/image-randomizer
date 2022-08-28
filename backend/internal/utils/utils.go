package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID to generate random uuid.
func GenerateUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

const passwordPepper string = "yagoobestgirl"

// EncodePassword to encode password with salt and pepper.
func EncodePassword(str, salt string) string {
	hasher := sha1.New()
	hasher.Write([]byte(salt + str + passwordPepper))
	return hex.EncodeToString(hasher.Sum(nil))
}
