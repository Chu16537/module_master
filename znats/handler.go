package znats

import (
	"github.com/Chu16537/gomodule/zjson"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

// 推送
func (h *Handler) Pub(msgAckID string, topicName string, msg interface{}) error {
	data, err := zjson.Marshal(msg)
	if err != nil {
		return errors.Wrapf(err, "Pub Marshal err :%v", err.Error())
	}

	if msgAckID == "" {
		if err := h.nc.Publish(topicName, data); err != nil {
			return errors.Wrapf(err, "nc Publish err :%v", err.Error())
		}
	} else {
		if _, err := h.js.Publish(h.ctx, topicName, data); err != nil {
			return errors.Wrapf(err, "js Publish err :%v", err.Error())
		}

		v, isLoad := h.msgAckMap.Load(msgAckID)
		if !isLoad {
			return nil
		}

		msg, ok := v.(jetstream.Msg)
		if !ok {
			return nil
		}

		msg.Ack()

		h.msgAckMap.Delete(msgAckID)
	}

	return nil
}

// 訂閱
func (h *Handler) Sub(streamName string, topicname string, f func(string, []byte)) error {
	var err error

	// 不需要stream
	if streamName == "" {
		if _, err = h.nc.Subscribe(topicname, func(m *nats.Msg) {
			if f != nil {
				f("", m.Data)
			}
		}); err != nil {
			return errors.Wrapf(err, "Subscribe err :%v", err.Error())
		}
		return nil
	}

	// 更新
	if err := h.delStreamTopic(streamName, topicname); err != nil {
		return err
	}

	// 取得 stream
	s, err := h.createStream(streamName, topicname)
	if err != nil {
		return err
	}

	//開始執行
	err = h.startSub(s, streamName, topicname, f)
	if err != nil {
		return err
	}

	return nil
}

// 取消訂閱
func (h *Handler) UnSub(streamName string, topicname string) {
	h.js.DeleteConsumer(h.ctx, streamName, topicname)
}
