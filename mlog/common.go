package mlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mtime"
	"github.com/sirupsen/logrus"
)

type LogData struct {
	Level  logrus.Level
	Tracer string // 會重複 追蹤使用 可能有一個請求 執行多個服務/多個模組
	FnName string
	Data   interface{}
	Err    *errorcode.Error
}

// UTCFormatter 自訂的 JSON 格式化器，強制使用 UTC 時間
type logOpt struct {
	logrus.JSONFormatter
}

func (o *logOpt) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Time = entry.Time.UTC()

	if _, ok := entry.Data["code"]; ok {
		// 格式化錯誤消息
		entry.Message = formatErrorMessage(entry.Message)
	}

	return o.JSONFormatter.Format(entry)
}

func formatErrorMessage(errMsg string) string {
	var sb strings.Builder
	lines := strings.Split(errMsg, "\n")
	for _, line := range lines {
		// 使用正則表達式去除多餘的空格和制表符
		re := regexp.MustCompile(`^\s+`)
		formattedLine := re.ReplaceAllString(line, "")
		sb.WriteString(" " + formattedLine)
	}
	return sb.String()
}

// 創建新的log file
func (h *handler) createNewFile() {
	if h.config.LogFilePath == "" {
		return
	}

	nowFormat := mtime.GetTimeFormatTime(h.t, mtime.Format_YMD)

	// 假如不同日期
	if h.currentDate != nowFormat {

		fp := filepath.Join(h.config.LogFilePath, fmt.Sprintf("%s_%s.log", h.config.LogFileName, nowFormat))

		// 確保目錄存在
		err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"topic": "mlog",
			}).Error("createNewFile MkdirAll err date:", h.currentDate)
			return
		}

		file, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"topic": "mlog",
			}).Error("createNewFile OpenFile err date:", h.currentDate)
			return
		}

		err = h.file.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"topic": "mlog",
			}).Error("createNewFile OpenFile err date:", h.currentDate)
			return
		}

		// 日期替換
		h.currentDate = nowFormat
		h.file = file

		logrus.SetOutput(h.file)
	}
}
