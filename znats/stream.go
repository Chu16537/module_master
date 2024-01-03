package znats

import (
	"github.com/Chu16537/gomodule/utils"
	"github.com/pkg/errors"

	"github.com/nats-io/nats.go/jetstream"
)

// 移除 stream的 topic
func (h *Handler) delStreamTopic(streamName, topicname string) error {
	stream, _ := h.js.StreamNameBySubject(h.ctx, topicname)

	// 現在要使用的 stream 與之前的不相同
	if stream != "" && stream != streamName {
		oldStream, err := h.js.Stream(h.ctx, stream)
		if err != nil {
			return errors.Wrapf(err, "delStreamTopic Stream err :%v", err.Error())
		}
		// stream資料
		info, _ := oldStream.Info(h.ctx)
		//是否刪除
		isDel := false
		for _, v := range info.Config.Subjects {
			if v == topicname {
				isDel = true
				break
			}
		}
		if isDel {
			newSubName := utils.RemoveItems(info.Config.Subjects, topicname)
			info.Config.Subjects = newSubName
			if _, err := h.js.UpdateStream(h.ctx, info.Config); err != nil {
				return errors.Wrapf(err, "delStreamTopic UpdateStream err :%v", err.Error())
			}
		}
	}

	return nil
}

// 產生新的 stream
func (h *Handler) createStream(streamName, topicname string) (jetstream.Stream, error) {
	var (
		s   jetstream.Stream
		err error
	)
	// 取得 stream
	s, err = h.js.Stream(h.ctx, streamName)
	if err != nil {
		return nil, errors.Wrapf(err, "createStream Stream err :%v", err.Error())
	}

	if s == nil {
		conf := jetstream.StreamConfig{
			Name:      streamName,
			Retention: jetstream.WorkQueuePolicy, //有收到ack 就刪除
			Subjects:  []string{topicname},
		}
		// 創建 stream
		s, err := h.js.CreateStream(h.ctx, conf)
		if err != nil {
			return nil, errors.Wrapf(err, "createStream CreateStream err :%v", err.Error())
		}
		return s, nil
	}

	// 判斷是否已經有topicname
	isSub := false
	for _, v := range s.CachedInfo().Config.Subjects {
		if v == topicname {
			isSub = true
			break
		}
	}

	// 假如沒有註冊topic
	if !isSub {
		// stream資料
		info, _ := s.Info(h.ctx)
		newSubName := append(info.Config.Subjects, topicname)
		info.Config.Subjects = newSubName
		s, err = h.js.UpdateStream(h.ctx, info.Config)
		if err != nil {
			return nil, errors.Wrapf(err, "createStream UpdateStream err :%v", err.Error())
		}
	}

	return s, nil
}

func (h *Handler) startSub(s jetstream.Stream, streamName, topicname string, f func([]byte)) error {
	c, err := s.CreateOrUpdateConsumer(h.ctx, jetstream.ConsumerConfig{
		Durable:       streamName, // 使用永久的
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: topicname, // 指定sub
	})

	if err != nil {
		return errors.Wrapf(err, "startSub CreateOrUpdateConsumer err :%v", err.Error())
	}

	go func() {
		consContext, err := c.Consume(func(msg jetstream.Msg) {
			metadata, _ := msg.Metadata()
			if metadata != nil {
				// 執行事件
				f(msg.Data())
				msg.Ack()
			}
		})

		if err != nil {
			h.js.DeleteConsumer(h.ctx, streamName, topicname)
		}

		defer consContext.Stop()
	}()

	return nil
}
