package mnats

import (
	"time"

	"github.com/nats-io/nats.go"
)

// 產生新的 stream
func (h *Handler) createStream(streamName string, liveSecond time.Duration, maxLen int64) error {
	// 取得 stream
	_, err := h.js.StreamInfo(streamName)
	if err != nil && err.Error() == "nats: stream not found" {
		// 沒有stream ，創建stream
		_, err = h.js.AddStream(&nats.StreamConfig{
			Name:    streamName,
			MaxAge:  liveSecond,
			MaxMsgs: maxLen, // 設置最大消息數量
			// Retention: nats.WorkQueuePolicy, // 僅保留未被消費的消息
		})
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// 創建 Subjects
func (h *Handler) AddSubjects(stream string, subjectName string) error {
	// 取得所有subject
	streamInfo, err := h.js.StreamInfo(stream)
	if err != nil {
		return err
	}

	// 更新
	_, err = h.js.UpdateStream(
		&nats.StreamConfig{
			Name:     stream,
			Subjects: append(streamInfo.Config.Subjects, subjectName),
		},
	)

	// 錯誤不是 已經創建過 subject
	if err != nil && err.Error() != "nats: duplicate subjects detected" {
		return err
	}

	return nil
}
