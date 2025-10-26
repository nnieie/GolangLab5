package service

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/nnieie/golanglab5/cmd/user/dal/cache"
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

const (
	qrcodeWidth = 256
	qrcodeHight = 256
)

func (s *UserService) GetMFAqrcode(userID int64) (string, string, error) {
	user, err := db.QueryUserByID(s.ctx, userID)
	if err != nil {
		return "", "", err
	}
	key, err := utils.GenerateTotp(user.UserName)
	if err != nil {
		return "", "", err
	}
	secret := key.Secret()
	qrcode, err := key.Image(qrcodeWidth, qrcodeHight)
	if err != nil {
		logger.Errorf("mfa qrcode generate err: %v", err)
		return "", "", err
	}

	// 将 qrcode 编码为 base64
	var buf bytes.Buffer
	err = png.Encode(&buf, qrcode)
	if err != nil {
		logger.Errorf("mfa qrcode encode err: %v", err)
		return "", "", err
	}
	base64Qrcode := base64.StdEncoding.EncodeToString(buf.Bytes())

	// 加密 secret
	encryptedSecret, err := utils.Encrypt(secret)
	if err != nil {
		return "", "", err
	}

	// 将 encryptedSecret 存入redis, 便于绑定验证
	err = cache.SetTOTPSecret(s.ctx, encryptedSecret, userID)
	if err != nil {
		return "", "", err
	}

	return secret, base64Qrcode, nil
}
