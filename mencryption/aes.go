package mencryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	AES_CBC = "CBC"
	AES_CTR = "CTR"
)

// Decode 根據 aesType 選擇不同的解密類型
func AesDecode(aesType string, key, iv []byte, ciphertext string, data any) error {
	iv = checkIV(iv)

	err := validateKeyIV(key, iv)
	if err != nil {
		return err
	}

	switch aesType {
	case AES_CBC:
		return decodeCBC(key, iv, ciphertext, data) // 使用 AES-CBC 解密
	case AES_CTR:
		return decodeCRT(key, iv, ciphertext, data) // 使用 AES-CTR 解密
	default:
		return fmt.Errorf("unsupported aesType: %s", aesType) // 不支援的 AES 類型
	}
}

// Encode 根據 aesType 選擇不同的加密類型
func AesEncode(aesType string, key, iv []byte, data any) ([]byte, error) {
	iv = checkIV(iv)

	err := validateKeyIV(key, iv)
	if err != nil {
		return nil, err
	}

	// json 編碼
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch aesType {
	case AES_CBC:
		return encodeCBC(key, iv, jsonData) // 使用 AES-CBC 加密
	case AES_CTR:
		return encodeCRT(key, iv, jsonData) // 使用 AES-CTR 加密
	default:
		return nil, fmt.Errorf("unsupported aesType: %s", aesType) // 不支援的 AES 類型
	}
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

func checkIV(iv []byte) []byte {
	if len(iv) == 0 {
		// 如果 IV 為空，則初始化為 16 字節全零
		return make([]byte, aes.BlockSize)
	}
	return iv
}

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

	// 更高效的填充檢查
	if !bytes.Equal(data[length-unpadding:], bytes.Repeat([]byte{byte(unpadding)}, unpadding)) {
		return nil, fmt.Errorf("pkcs7Unpadding: invalid padding")
	}

	return data[:(length - unpadding)], nil
}

// AesDecrypt AES CBC 解密
func decodeCBC(key, iv []byte, ciphertext string, data any) error {
	encryptedBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	plaintext := make([]byte, len(encryptedBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encryptedBytes)

	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return err
	}

	// json 解碼
	if err := json.Unmarshal(plaintext, data); err != nil {
		return err
	}

	return nil
}

// AesEncrypt AES CBC 加密
func encodeCBC(key, iv []byte, jsonData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	jsonData = pkcs7Padding(jsonData, aes.BlockSize)
	ciphertext := make([]byte, len(jsonData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, jsonData)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// AES-CTR 解密 Base64 數據並解析 JSON
func decodeCRT(key, iv []byte, ciphertext string, data any) error {
	// Base64 解碼
	encryptedBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return err
	}

	// AES 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(encryptedBytes, encryptedBytes)

	// JSON 解析
	err = json.Unmarshal(encryptedBytes, data)
	if err != nil {
		return err
	}
	return nil
}

// AES-CTR 加密數據並返回 Base64 編碼的結果
func encodeCRT(key, iv []byte, jsonData []byte) ([]byte, error) {
	// aes 加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(jsonData, jsonData)

	return []byte(base64.StdEncoding.EncodeToString(jsonData)), nil
}
