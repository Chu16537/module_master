package mconfig

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func LoadEnv(path string) error {
	return godotenv.Load(path)
}

// 取得字串
func GetEnvString(key string) string {
	return os.Getenv(key)
}

// 取得字串 請使用逗號"," 做區隔
func GetEnvStringArray(key string) []string {
	s := os.Getenv(key)
	// 將逗號分隔的字符串轉換為 []string
	return strings.Split(s, ",")
}

// 取得數字
func GetEnvDecimal(key string) (decimal.Decimal, error) {
	return decimal.NewFromString(os.Getenv(key))
}

func GetFloat64(key string) (float64, error) {
	d, err := decimal.NewFromString(os.Getenv(key))
	if err != nil {
		return 0, err
	}
	return d.InexactFloat64(), nil
}

func GetInt64(key string) (int64, error) {
	d, err := decimal.NewFromString(os.Getenv(key))
	if err != nil {
		return 0, err
	}
	return d.IntPart(), nil
}

// 取得時間
func GetDuration(key string) (time.Duration, error) {
	return time.ParseDuration(os.Getenv(key))
}
