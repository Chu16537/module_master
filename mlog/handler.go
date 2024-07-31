package mlog

import (
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
)

// 不上傳到遠端
func (h *Handler) Debug(fnName string, msg interface{}) {
	l := Log{
		LV:       "debug",
		FuncName: fnName,
		Msg:      msg,
	}

	fmt.Println(l)
}

// 上傳到遠端
func (h *Handler) Info(fnName string, msg interface{}) {
	l := Log{
		LV:       "info",
		FuncName: fnName,
		Msg:      msg,
	}

	fmt.Println(l)
}

// 可忽略的錯誤 上傳到遠端
func (h *Handler) Warn(fnName string, msg *errorcode.Error) {
	l := Log{
		LV:       "warn",
		FuncName: fnName,
		Code:     msg.Code(),
		Msg:      fmt.Sprintf("%+v", msg.Err()),
	}

	fmt.Println(l)
}

// 錯誤 上傳到遠端
func (h *Handler) Error(fnName string, msg *errorcode.Error) {
	l := Log{
		LV:       "error",
		FuncName: fnName,
		Code:     msg.Code(),
		Msg:      fmt.Sprintf("%+v", msg.Err()),
	}

	fmt.Println(l)
}
