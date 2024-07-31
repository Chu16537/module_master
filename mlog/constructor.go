package mlog

import "github.com/Chu16537/module_master/errorcode"

type ILog interface {
	Debug(fnName string, msg interface{})
	Info(fnName string, msg interface{})
	Warn(funcname string, msg *errorcode.Error)
	Error(funcname string, msg *errorcode.Error)
}

type Log struct {
	LV       string
	FuncName string
	Code     int
	Msg      interface{}
}

type Handler struct {
	Addr string
}

func New(addr string) ILog {
	return &Handler{
		Addr: addr,
	}
}
