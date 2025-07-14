package utils

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateMFASecret genera un nuevo secreto TOTP para un usuario
func GenerateMFASecret(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Menchaca System",
		AccountName: email,
		SecretSize:  20,
	})
}

// ValidateTOTP verifica un código TOTP
func ValidateTOTP(token string, secret string) bool {
	return totp.Validate(token, secret)
}

