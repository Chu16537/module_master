package muid

import "github.com/rs/xid"

// 取得唯一碼
func GetUID() string {
	return xid.New().String()
}
