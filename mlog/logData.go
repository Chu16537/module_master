package mlog

import (
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/muid"
	"github.com/sirupsen/logrus"
)

type LogData struct {
	Level  logrus.Level
	Server string
	Tracer string // 會重複 追蹤使用 可能有一個請求 執行多個服務/多個模組
	FnName string
	Data   interface{}
	Err    *errorcode.Error
}

func NewLogData(fnName string, nodeID int64) *LogData {
	l := new(LogData)
	l.FnName = fnName
	l.Tracer = fmt.Sprintf("%v_%v", muid.CreatRandomString(10), nodeID)

	return l
}
