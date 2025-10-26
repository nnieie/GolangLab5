package utils

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTotp 生成 TOTP 密钥
func GenerateTotp(userName string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "tkrpc",
		AccountName: userName,
	})
}

// CheckTotp 验证 TOTP
func CheckTotp(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
