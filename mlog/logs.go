package mlog

import (
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/muid"
	"github.com/sirupsen/logrus"
)

type ILog interface {
	Debug(opt *LogData)
	Info(opt *LogData)
	Warn(opt *LogData)
	Error(opt *LogData)
}

type LogData struct {
	Level  logrus.Level
	Server string
	Tracer string // 會重複 追蹤使用 可能有一個請求 執行多個服務/多個模組
	FnName string
	Data   interface{}
	Err    *errorcode.Error
}

func NewLogData(fnName string) *LogData {
	l := new(LogData)
	l.FnName = fnName
	l.Tracer = fmt.Sprintf("%v_%v", muid.CreatRandomString(10), h.nodeId)

	return l
}

func (l *Log) initial(opt *LogData) logrus.Fields {
	l.handler.createNewFile()

	f := logrus.Fields{}
	f["server"] = serverName

	if opt == nil {
		return f
	}

	if opt.Tracer != "" {
		f["tracer"] = opt.Tracer
	}

	if opt.FnName != "" {
		f["topic"] = fmt.Sprintf("%v_%v", l.name, opt.FnName)
	} else {
		f["topic"] = l.name
	}

	if opt.Err != nil {
		f["code"] = opt.Err.Code()
	}

	return f
}

func (l *Log) Debug(opt *LogData) {
	fields := l.initial(opt)
	if opt.Data != nil {
		logrus.WithFields(fields).Debug(opt.Data)
	} else {
		logrus.WithFields(fields).Debug("")
	}
}

func (l *Log) Info(opt *LogData) {
	fields := l.initial(opt)
	if opt.Data != nil {
		logrus.WithFields(fields).Info(opt.Data)
	} else {
		logrus.WithFields(fields).Info("")
	}
}

func (l *Log) Warn(opt *LogData) {
	fields := l.initial(opt)
	if opt.Err.Err() != nil {
		logrus.WithFields(fields).Warn(fmt.Sprintf("%+v", opt.Err.Err()))
	} else {
		logrus.WithFields(fields).Warn("")
	}
}

func (l *Log) Error(opt *LogData) {
	fields := l.initial(opt)
	if opt.Err.Err() != nil {
		logrus.WithFields(fields).Error(fmt.Sprintf("%+v", opt.Err.Err()))
	} else {
		logrus.WithFields(fields).Error("")
	}
}
