package mencryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// pkcs7Padding PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpadding 移除 PKCS7 填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("pkcs7Unpadding: input data is empty")
	}

	unpadding := int(data[length-1])
	if unpadding < 1 || unpadding > aes.BlockSize || unpadding > length {
		return nil, fmt.Errorf("pkcs7Unpadding: invalid padding size")
	}

	// More efficient padding check
	if !bytes.Equal(data[length-unpadding:], bytes.Repeat([]byte{byte(unpadding)}, unpadding)) {
		return nil, fmt.Errorf("pkcs7Unpadding: invalid padding")
	}

	return data[:(length - unpadding)], nil
}

// AesEncrypt AES CBC 加密
func AesEncrypt(plaintext, key, iv []byte) (string, error) {
	if err := validateKeyIV(key, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("AesEncrypt: failed to create cipher: %w", err)
	}

	plaintext = pkcs7Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesDecrypt AES CBC 解密
func AesDecrypt(ciphertext string, key, iv []byte) ([]byte, error) {
	if err := validateKeyIV(key, iv); err != nil {
		return nil, err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("AesDecrypt: failed to decode base64: %w", err)
	}

	if len(encryptedBytes)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("AesDecrypt: encryptedBytes is not a multiple of block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("AesDecrypt: failed to create cipher: %w", err)
	}

	plaintext := make([]byte, len(encryptedBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encryptedBytes)

	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return nil, fmt.Errorf("AesDecrypt: %w", err)
	}

	return plaintext, nil
}

func validateKeyIV(key, iv []byte) error {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("invalid key length %d, must be 16, 24, or 32 bytes", keyLen)
	}
	if len(iv) != aes.BlockSize {
		return fmt.Errorf("invalid IV length %d, must be %d bytes", len(iv), aes.BlockSize)
	}
	return nil
}
