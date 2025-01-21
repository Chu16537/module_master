package mencryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// Padding: PKCS7 填充方式
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// Unpadding: 移除 PKCS7 填充
func pkcs7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// AES 加密
func AesEncrypt(plaintext, key []byte, iv []byte) (string, error) {
	if iv == nil || key == nil {
		return "", fmt.Errorf("AesEncrypt key:%v iv:%v nil", key, iv)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = pkcs7Padding(plaintext, block.BlockSize())

	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 將加密數據轉換為 Base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	return encryptedBase64, nil
}

// AES 解密
func AesDecrypt(ciphertext string, key []byte, iv []byte) ([]byte, error) {
	if iv == nil || key == nil {
		return nil, fmt.Errorf("AesDecrypt key:%v iv:%v nil", key, iv)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密：從 Base64 還原密文
	encryptedBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	origData := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(origData, encryptedBytes)

	plaintext := pkcs7Unpadding(origData)

	return plaintext, nil
}
