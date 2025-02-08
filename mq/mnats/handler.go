package mnats

import (
	"github.com/Chu16537/module_master/mq"
	"github.com/Chu16537/module_master/proto"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// 發布消息
func (h *Handler) Pub(subject string, data []byte) error {
	if err := h.nc.Publish(subject, data); err != nil {
		return err
	}
	return nil
}

// 訂閱
func (h *Handler) Sub(subject string, consumer string, mode mq.SubMode, subChan chan proto.MQSubData) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.js == nil {
		return errors.New("nats nil")
	}

	opts := []nats.SubOpt{}

	switch mode.Mode {
	case mq.Sub_Mode_Last_Ack:
		opts = append(opts, nats.Durable(consumer)) // 從最後ack 開始
	case mq.Sub_Mode_Last:
		opts = append(opts, nats.DeliverNew()) // 從最新消息開始
	case mq.Sub_Mode_Sequence:
		if mode.StartSequenceID < 1 {
			mode.StartSequenceID = 1
		}
		opts = append(opts, nats.StartSequence(mode.StartSequenceID)) // 從指定的 Sequence 開始
	}

	// 訂閱
	sub, err := h.js.Subscribe(subject, func(msg *nats.Msg) {
		data, _ := msg.Metadata()
		subChan <- proto.MQSubData{
			Data:       msg.Data,
			SequenceID: data.Sequence.Consumer,
		}

		// 要記錄最後ack
		if mode.Mode == mq.Sub_Mode_Last_Ack {
			msg.Ack()
		}
	}, opts...)

	if err != nil {
		return errors.Wrap(err, "訂閱失敗")
	}

	// 記錄訂閱
	h.subMap[subject] = sub
	return nil
}

// 取消訂閱
func (h *Handler) UnSub(consumer string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if sub, ok := h.subMap[consumer]; ok {
		sub.Unsubscribe()
		delete(h.subMap, consumer)
	}
}

// 刪除指定 subject 內的所有訊息
func (h *Handler) DelSubject(streamName, subjectName string) error {
	if h.js == nil {
		return errors.New("NATS JetStream 尚未初始化")
	}

	err := h.js.PurgeStream(streamName, &nats.StreamPurgeRequest{Subject: subjectName})
	if err != nil {
		return errors.Wrap(err, "清除 subject 訊息失敗")
	}

	return nil
}

// 刪除 Stream
func (h *Handler) DelStream(streamName string) error {
	if h.js == nil {
		return errors.New("nats nil")
	}

	err := h.js.DeleteStream(streamName)
	if err != nil {
		return errors.Wrap(err, "刪除 Stream 失敗")
	}
	return nil
}
