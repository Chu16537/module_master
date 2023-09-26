package znats

import (
	"fmt"

	"github.com/Chu16537/gomodule/utils"

	"github.com/nats-io/nats.go/jetstream"
)

// 移除 stream的 topic
func (h *Handler) delStreamTopic(streamName, topicname string) {
	stream, _ := h.js.StreamNameBySubject(h.ctx, topicname)

	// 現在要使用的 stream 與之前的不相同
	if stream != "" && stream != streamName {
		oldStream, err := h.js.Stream(h.ctx, stream)
		if err != nil {
			fmt.Println("delStreamTopic", err)
			return
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
			_, err := h.js.UpdateStream(h.ctx, info.Config)
			if err != nil {
				fmt.Println("delStreamTopic UpdateStream", err)
				return
			}
		}
	}
}

// 產生新的 stream
func (h *Handler) createStream(streamName, topicname string) jetstream.Stream {
	// 取得 stream
	s, _ := h.js.Stream(h.ctx, streamName)

	if s == nil {
		conf := jetstream.StreamConfig{
			Name:      streamName,
			Retention: jetstream.WorkQueuePolicy, //有收到ack 就刪除
			Subjects:  []string{topicname},
		}
		// 創建 stream
		s, err := h.js.CreateStream(h.ctx, conf)
		if err != nil {
			fmt.Println("createStream", err)
			return nil
		}
		return s
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
		s, _ = h.js.UpdateStream(h.ctx, info.Config)
	}

	return s
}
