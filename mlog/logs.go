package mlog

import (
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/sirupsen/logrus"
)

type ILog interface {
	New(level int, fnName string, tracer string, data interface{}, err *errorcode.Error)
}
type LogData struct {
	Level  int              // 日誌級別
	Server string           // 伺服器名稱
	Tracer string           // 追蹤 ID，便於在分散式系統中追蹤多個日誌請求
	FnName string           // 方法名稱或主題，標識日誌記錄所屬的功能模組。
	Data   interface{}      // 日誌中包含的主要資料
	Err    *errorcode.Error // 錯誤信息，使用自定義的 errorcode
}

func (l *Log) New(level int, fnName string, tracer string, data interface{}, err *errorcode.Error) {

	l.handler.createNewFile()
	fields := logrus.Fields{
		"server": l.handler.config.Name,
		"tracer": tracer,
		"topic":  fmt.Sprintf("%v_%v", l.name, fnName),
	}

	if err != nil {
		fields["code"] = err.GetCode()
	}

	entry := logrus.WithFields(fields)

	msg := l.formatMessage(level, data, err)

	switch level {
	case DebugLevel:
		entry.Debug(msg)
	case InfoLevel:
		entry.Info(msg)
	case WarnLevel:
		entry.Warn(msg)
	case ErrorLevel:
		entry.Error(msg)
	}

}

func (l *Log) formatMessage(level int, data interface{}, err *errorcode.Error) interface{} {
	if level == WarnLevel || level == ErrorLevel {
		if err != nil && err.GetErr() != nil {
			return fmt.Sprintf("%+v", err.GetErr())
		}
	} else {
		if data != nil {
			return data
		}
	}

	return ""
}
