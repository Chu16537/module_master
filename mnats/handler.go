package mnats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

// 推送
func (h *Handler) Pub(topicName string, data []byte) error {
	if err := h.nc.Publish(topicName, data); err != nil {
		return errors.Wrapf(err, "nc Publish err :%v", err.Error())
	}

	return nil
}

// 推送
func (h *Handler) PubStream(topicName string, data []byte) error {
	if _, err := h.js.Publish(h.ctx, topicName, data); err != nil {
		return errors.Wrapf(err, "js Publish err :%v", err.Error())
	}
	return nil
}

// 訂閱
func (h *Handler) Sub(topicname string, f func([]byte)) error {
	if _, err := h.nc.Subscribe(topicname, func(msg *nats.Msg) {
		f(msg.Data)
	}); err != nil {
		return errors.Wrapf(err, "Subscribe err :%v", err.Error())
	}
	return nil
}

// 訂閱
func (h *Handler) SubStream(streamName string, topicname string, f func([]byte)) error {
	var err error

	// 更新
	err = h.delStreamTopic(streamName, topicname)
	if err != nil {
		return err
	}

	// 取得 stream
	err = h.createStream(streamName, topicname)
	if err != nil {
		return err
	}

	return nil
}

// 取得sub msg
func (h *Handler) GetMsg(streamName string, topicName string, f func([]byte), count int) {
	v, isLoad := h.consumerMap.Load(topicName)

	if !isLoad {
		h.SubStream(streamName, topicName, f)
		v, isLoad = h.consumerMap.Load(topicName)
		if !isLoad {
			return
		}
	}

	c, ok := v.(jetstream.Consumer)

	if !ok {
		return
	}

	msgBatch, err := c.Fetch(count)
	if err != nil {
		return
	}

	msgs := msgBatch.Messages()
	msgLen := len(msgs)
	idx := 0
	for msg := range msgs {
		//消費資料
		f(msg.Data())
		idx++

		//最後一筆資料 通知ack
		if idx == msgLen-1 {
			msg.Ack()
		}
	}
}

// 取消訂閱
func (h *Handler) UnSub(streamName string, topicname string) {
	h.js.DeleteConsumer(h.ctx, streamName, topicname)

	h.consumerMap.Delete(topicname)
}
