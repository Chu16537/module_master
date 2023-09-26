package znats

import (
	"fmt"
	"time"

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

func (h *Handler) startSub(c jetstream.Consumer, f func(interface{})) {
	info, _ := c.Info(h.ctx)
	fmt.Println("StartSub", info.Config.FilterSubject)

	cons, _ := c.Consume(func(msg jetstream.Msg) {
		data := string(msg.Data())
		fmt.Println(info.Config.FilterSubject, "time", time.Now())
		if f != nil {
			// 執行事件結束
			f(data)
		}
		// 回傳ack
		msg.Ack()
		fmt.Println(info.Config.FilterSubject, "time", time.Now())
	})
	h.consumeMap[c.CachedInfo().Config.FilterSubject] = cons
}
