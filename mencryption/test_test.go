package mencryption_test

import (
	"fmt"
	"testing"

	"github.com/Chu16537/module_master/mencryption"
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

	ty := mencryption.AES_CBC
	// ty := mencryption.AES_CBC

	aesEncStr, err := mencryption.AesEncode(ty, []byte(key), []byte(iv), a)
	if err != nil {
		t.Error("AesEncode", err)
		return
	}

	fmt.Println("aesEncStr", string(aesEncStr))

	aa := &A{}
	err = mencryption.AesDecode(ty, []byte(key), []byte(iv), string(aesEncStr), aa)
	if err != nil {
		t.Errorf("AesDecode %+v", err)
		return
	}

	fmt.Println("aa", aa)

}
