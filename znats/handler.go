package znats

import (
	"github.com/nats-io/nats.go"
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
	var err error

	if _, err = h.nc.Subscribe(topicname, func(m *nats.Msg) {
		f(m.Data)
	}); err != nil {
		return errors.Wrapf(err, "Subscribe err :%v", err.Error())
	}
	return nil
}

// 訂閱
func (h *Handler) SubStream(streamName string, topicname string, f func(uint64, []byte)) error {
	var err error

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
