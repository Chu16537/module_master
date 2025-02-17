package mencryption_test

import (
	"fmt"
	"testing"

	"github.com/Chu16537/module_master/mencryption"
	"github.com/Chu16537/module_master/mjson"
)

func TestAes(t *testing.T) {
	key := "qqqqwwwweeeerrrr"
	iv := "0000000000000000"

	type A struct {
		Account  string `json:"account"`
		GameId   string `json:"game_id"`
		Currency string `json:"currency"`
		Lang     string `json:"lang"`
	}

	a := &A{
		Account:  "test_Account",
		GameId:   "1",
		Currency: "NTD",
		Lang:     "en",
	}

	b, err := mjson.Marshal(a)
	if err != nil {
		t.Error("Marshal", err)
		return
	}

	aesEncStr, err := mencryption.AesEncrypt(b, []byte(key), []byte(iv))
	if err != nil {
		t.Error("AesEncrypt", err)
		return
	}

	fmt.Println("aesEncStr", string(aesEncStr))

	aesDecStr, err := mencryption.AesDecrypt(aesEncStr, []byte(key), []byte(iv))
	if err != nil {
		t.Error("AesDecrypt", err)
		return
	}

	fmt.Println("aesDecStr", string(aesDecStr))

	s := &A{}

	err = mjson.Unmarshal(aesDecStr, s)
	if err != nil {
		t.Error("Unmarshal", err)
		return
	}

	fmt.Println("s", s)
}
