package mlog

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	Name string `env:"LOG_NAME" yaml:"name"` // 服務名稱
}

type handler struct {
	config *Config
	log    zerolog.Logger
}

var h *handler

func Init(config *Config) error {
	if config.Name == "" {
		return errors.New("config name is nil")
	}

	// 设置全局时间格式为 Unix 时间戳
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zlog := zerolog.New(os.Stdout). // 输出到标准输出
					With().
					Timestamp().                     // 为每条日志添加时间戳
					Str("server-name", config.Name). // 添加自定义字段
					Logger().
					Level(zerolog.DebugLevel) // 设置日志等级为 Debug

	h = &handler{
		config: config,
		log:    zlog,
	}

	return nil
}

func Info(args ...interface{}) {
	h.log.WithLevel(zerolog.InfoLevel).CallerSkipFrame(1).Msg(fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.InfoLevel).Msgf(format, args...)
}

func Trace(args ...interface{}) {
	h.log.WithLevel(zerolog.TraceLevel).Msg(fmt.Sprint(args...))
}

func Tracef(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.TraceLevel).Msgf(format, args...)
}

func Debug(args ...interface{}) {
	h.log.WithLevel(zerolog.DebugLevel).Msg(fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.DebugLevel).Msgf(format, args...)
}

func Warn(args ...interface{}) {
	h.log.WithLevel(zerolog.WarnLevel).Caller(1).Msg(fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.WarnLevel).Caller(1).Msgf(format, args...)
}

func Error(args ...interface{}) {
	h.log.WithLevel(zerolog.ErrorLevel).Caller(1).Msg(fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.ErrorLevel).Caller(1).Msgf(format, args...)
}

func Fatal(args ...interface{}) {
	h.log.WithLevel(zerolog.FatalLevel).Caller(3).Msg(fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	h.log.WithLevel(zerolog.FatalLevel).Caller(3).Msgf(format, args...)
	os.Exit(1)
}
