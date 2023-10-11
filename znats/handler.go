package znats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// 推送
func (h *Handler) Pub(topicName string, msg interface{}) {
	payload := []byte(fmt.Sprintln(msg))
	h.js.Publish(h.ctx, topicName, payload)
}

// 訂閱
func (h *Handler) Sub(streamName string, topicname string, f func([]byte)) bool {
	var err error

	// 不需要stream
	if streamName == "" {
		_, err = h.nc.Subscribe(topicname, func(m *nats.Msg) {
			if f != nil {
				f(m.Data)
			}
		})

		return err == nil
	}

	// 更新
	h.delStreamTopic(streamName, topicname)

	// 取得 stream
	s := h.createStream(streamName, topicname)

	if s == nil {
		fmt.Println("stream nill")
		return false
	}

	c := h.getCons(s, topicname)

	if c == nil {
		fmt.Println("cons nill")
		return false
	}

	//開始執行
	go h.startSub(c, f)

	return true
}

// 取消訂閱
func (h *Handler) UnSub(topicname string) {
	h.consumeMap[topicname].Stop()
}
