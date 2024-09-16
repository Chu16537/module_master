package muid

import (
	"fmt"

	"github.com/rs/xid"
)

// 取得唯一碼
func GetUID() string {
	return xid.New().String()
}

// 取得唯一碼 要node
func GetUIDNode(node int) string {
	return fmt.Sprintf("%v_%v", GetUID(), node)
}
