package mlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Chu16537/module_master/mtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name         string // 伺服器名稱，用於標識日誌來源
	FilePath     string // 本地日誌檔案路徑
	ElasticURL   string // Elasticsearch 伺服器 URL，用於將日誌輸出到 Elasticsearch
	ElasticIndex string // Elasticsearch 索引，用於分組管理日誌
}

type Log struct {
	handler *handler
	name    string
}
type handler struct {
	config      *Config   // 配置信息
	file        *os.File  // 當前日誌檔案指標
	t           time.Time // 當前日誌的日期時間
	currentDate string    // 當前日期
}

var (
	h          *handler
	serverName string
)

func New(config *Config) error {
	if config.Name == "" {
		errors.New("name is nil")
	}
	serverName = config.Name

	// 設定日誌格式
	opt := &logOpt{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	}

	logrus.SetFormatter(opt)
	logrus.SetLevel(logrus.DebugLevel)

	h = &handler{
		config:      config,
		t:           mtime.GetZero(),
		currentDate: mtime.GetTimeFormatTime(mtime.GetZero(), mtime.Format_YMD),
	}

	// 輸出到本地
	if config.FilePath != "" {
		fp := filepath.Join(config.FilePath, fmt.Sprintf("%s_%s.log", config.Name, h.currentDate))

		// 確保目錄存在
		err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if err != nil {
			return err
		}

		h.file, err = os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}

		logrus.SetOutput(h.file)
	} else {
		logrus.SetOutput(os.Stdout)
	}

	// 添加 Elasticsearch hook
	if config.ElasticURL != "" {
		esHook, err := newElasticsearchHook(config.ElasticURL, config.ElasticIndex)
		if err != nil {
			return err
		}

		logrus.AddHook(esHook)
	}

	return nil
}

func Get(name string) ILog {
	return &Log{
		handler: h,
		name:    name,
	}
}

func Done() {
	if h.file != nil {
		h.file.Close()
	}
}
