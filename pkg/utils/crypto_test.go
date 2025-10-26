package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHash(t *testing.T) {
	// 测试密码
	password := "114514"

	// 加密
	hashedPassword, err := Crypt(password)

	// 用 assert 来断言 加密过程不应返回 err
	assert.NoError(t, err, "crypt err")
	// 加密后不应为空
	assert.NotEmpty(t, hashedPassword, "crypt return empty")

	// 使用 t.Run 来创建子测试
	t.Run("verify correct password", func(t *testing.T) {
		isCorrect := VerifyPassword(password, hashedPassword)
		// 断言 isCorrect 应该是 true
		assert.True(t, isCorrect, "verify correct password should return true")
	})

	t.Run("verify wrong password", func(t *testing.T) {
		wrongPassword := "wrong"
		isCorrect := VerifyPassword(wrongPassword, hashedPassword)
		// 断言 isCorrect 应该是 false
		assert.False(t, isCorrect, "verify wrong password should return false")
	})
}

func TestEncryptAndDecrypt(t *testing.T) {
	// 测试用例
	cases := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "普通字符串",
			plaintext: "miao miao miao ~",
		},
		{
			name:      "汉字字符串",
			plaintext: "喵喵喵 ~",
		},
		{
			name:      "空字符串",
			plaintext: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// 加密
			ciphertext, err := Encrypt(tc.plaintext)
			assert.NoError(t, err, "encrypt err")

			// 解密
			decryptedText, err := Decrypt(ciphertext)
			assert.NoError(t, err, "decrypt err")

			// 对比
			assert.Equal(t, tc.plaintext, decryptedText, "decrypted text should match original plaintext")
		})
	}
}
