package mnats

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// 產生新的 Stream
func (h *Handler) createStream(streamName string, liveSecond time.Duration, maxLen int64) error {
	// 嘗試取得 Stream，如果不存在則創建
	if _, err := h.js.StreamInfo(streamName); err != nil {
		if !errors.Is(err, nats.ErrStreamNotFound) {
			return err
		}

		// 創建 Stream
		_, err = h.js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			MaxAge:    liveSecond,
			MaxMsgs:   maxLen,
			Retention: nats.WorkQueuePolicy, // 只有當訂閱者 ACK 之後才會刪除消息
		})

		if err != nil {
			return err
		}

	}

	return nil
}

// 創建 Subjects
func (h *Handler) AddSubjects(stream, subjectName string) error {
	streamInfo, err := h.js.StreamInfo(stream)
	if err != nil {
		return err
	}

	_, err = h.js.UpdateStream(&nats.StreamConfig{
		Name:      stream,
		Subjects:  append(streamInfo.Config.Subjects, subjectName),
		MaxAge:    streamInfo.Config.MaxAge,
		MaxMsgs:   streamInfo.Config.MaxMsgs,
		Retention: streamInfo.Config.Retention, // 只有當訂閱者 ACK 之後才會刪除消息
	})

	// 錯誤不是 已經創建過 subject
	if err != nil && err.Error() != "nats: duplicate subjects detected" {
		return err
	}
	return nil
}
