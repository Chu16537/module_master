package mjson

import (
	"bytes"

	"github.com/bytedance/sonic"
)

func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}

func NewEncoder(b *bytes.Buffer) sonic.Encoder {
	return sonic.ConfigDefault.NewEncoder(b)
}
