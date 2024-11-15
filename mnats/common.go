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
func (h *Handler) createSubjects(subjectName string) error {
	_, err := h.js.AddStream(
		&nats.StreamConfig{
			Name:     h.config.StreamName,
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

	// 已經創建過stream 就使用 UpdateStream
	_, err = h.js.UpdateStream(&nats.StreamConfig{
		Name:     h.config.StreamName,
		Subjects: []string{subjectName},
	})

	if err != nil {
		return err
	}

	// // 創建消費者
	// err = h.createConsumer(subjectName)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// // 創建 Consumer
// func (h *Handler) createConsumer(consumerName string) error {
// 	_, err := h.js.AddConsumer(h.config.StreamName, &nats.ConsumerConfig{
// 		Durable:   consumerName,
// 		AckPolicy: nats.AckAllPolicy,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
