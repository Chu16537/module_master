package mnats

import (
	"github.com/Chu16537/module_master/mmq"
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
func (h *Handler) Sub(subject string, consumer string, mode mmq.SubMode, subChan chan mmq.MQSubData) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.js == nil {
		return errors.New("nats nil")
	}

	var opts []nats.SubOpt

	switch mode.Mode {
	case mmq.Sub_Mode_Last_Ack:
		opts = append(opts, nats.Durable(consumer), nats.ManualAck(), nats.AckExplicit()) // 從最後ack 開始
	case mmq.Sub_Mode_Last:
		opts = append(opts, nats.DeliverNew()) // 從最新消息開始
	case mmq.Sub_Mode_Sequence:
		if mode.StartSequenceID < 1 {
			mode.StartSequenceID = 1
		}
		opts = append(opts, nats.StartSequence(mode.StartSequenceID)) // 從指定的 Sequence 開始
	}

	// 訂閱
	sub, err := h.js.Subscribe(subject, func(msg *nats.Msg) {
		data, _ := msg.Metadata()

		c := mmq.MQSubData{
			Data:       msg.Data,
			SequenceID: data.Sequence.Consumer,
			// Ack: func() {
			// 	msg.Ack()
			// },
		}

		subChan <- c

	}, opts...)

	if err != nil {
		return err
	}

	// 記錄訂閱
	h.subMap[subject] = sub
	return nil
}

// 取消訂閱
func (h *Handler) UnSub(consumer string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if sub, ok := h.subMap[consumer]; ok {
		sub.Unsubscribe()
		delete(h.subMap, consumer)
	}

	return nil
}

// 刪除指定 subject 內的所有訊息
func (h *Handler) DelSub(streamName, subjectName string) error {
	if h.js == nil {
		return errors.New("JetStream is not initialized")
	}

	err := h.js.PurgeStream(streamName, &nats.StreamPurgeRequest{Subject: subjectName})
	if err != nil {
		return err
	}

	return nil
}

// 刪除 Stream
func (h *Handler) DelStream(streamName string) error {
	if h.js == nil {
		return errors.New("JetStream is not initialized")
	}

	err := h.js.DeleteStream(streamName)
	if err != nil {
		return err
	}
	return nil
}
