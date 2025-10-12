package utils

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTotp(userName string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "tkrpc",
		AccountName: userName,
	})
}

func CheckTotp(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
