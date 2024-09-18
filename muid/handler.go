package muid

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/xid"
)

// 定義字母和數字的字元集
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 取得唯一碼
func GetUID() string {
	return xid.New().String()
}

// 取得唯一碼 要node
func GetUIDNode(node int) string {
	return fmt.Sprintf("%v_%v", GetUID(), node)
}

// 隨機字符串生成函數
func CreatRandomString(length int) string {
	// 預先生成隨機數的 buffer
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	buffer := make([]byte, 10) // 用來批量生成隨機數字節
	charsetLength := len(charset)

	for i := 0; i < length; i += 10 {
		// 批量生成隨機數
		seededRand.Read(buffer)

		// 將隨機數轉換為對應的字母和數字
		for j := 0; j < 10 && i+j < length; j++ {
			// 這裡使用 buffer[j] % charsetLength 來確保隨機數在字符集範圍內
			result[i+j] = charset[int(buffer[j])%charsetLength]
		}
	}
	return string(result)
}
