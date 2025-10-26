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

// CheckTotp 验证AES加密后的 TOTP
func CheckTotp(passcode string, encryptedSecret string) bool {
	secret, err := Decrypt(encryptedSecret)
	if err != nil {
		return false
	}
	return totp.Validate(passcode, secret)
}
