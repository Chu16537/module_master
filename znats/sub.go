package znats

import (
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (h *Handler) getCons(s jetstream.Stream, topicname string) jetstream.Consumer {
	c, err := s.CreateOrUpdateConsumer(h.ctx, jetstream.ConsumerConfig{
		Durable:       topicname, // 使用永久的
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: topicname, // 指定sub
	})

	if err != nil {
		fmt.Println("getCons", err)
		return nil
	}

	return c
}

func (h *Handler) startSub(c jetstream.Consumer, f func([]byte)) {
	cons, _ := c.Consume(func(msg jetstream.Msg) {
		if f != nil {
			// 執行事件
			f(msg.Data())
		}
		// 回傳ack
		msg.Ack()
	})
	h.consumeMap[c.CachedInfo().Config.FilterSubject] = cons
}
