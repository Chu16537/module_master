package mnats

import (
	"time"

	"github.com/nats-io/nats.go"
)

// 產生新的 stream
func (h *Handler) createStream(streamName string) error {
	// 取得 stream
	_, err := h.js.StreamInfo(streamName)
	if err != nil && err.Error() == "nats: stream not found" {
		// 沒有stream ，創建stream
		_, err = h.js.AddStream(&nats.StreamConfig{
			Name: streamName,
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
func (h *Handler) CreateSubjects(stream string, subjectName string) error {
	streamInfo, err := h.js.AddStream(
		&nats.StreamConfig{
			Name:     stream,
			Subjects: []string{subjectName},
			MaxAge:   time.Hour, // 設置消息存活時間為1小時
			// MaxMsgs:   1024, // 設置最大消息數量
			Retention: nats.WorkQueuePolicy, // 僅保留未被消費的消息
		},
	)

	// 錯誤不是 已經創建過stream
	if err != nil && err.Error() != "nats: stream name already in use" {
		return err
	}

	if streamInfo == nil {
		streamInfo, err = h.js.StreamInfo(stream)
		if err != nil {
			return err
		}
	}

	subjects := []string{}
	if streamInfo != nil {
		subjects = streamInfo.Config.Subjects
	}

	// 判斷是否有重複，假如有不更新
	for _, v := range subjects {
		if subjectName == v {
			return nil
		}
	}

	// 添加新的 Subject
	newSubjects := append(subjects, subjectName)

	// 已經創建過stream 就使用 UpdateStream
	_, err = h.js.UpdateStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: newSubjects,
	})

	if err != nil {
		return err
	}

	return nil
}
