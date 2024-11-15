package mnats

import (
	"fmt"

	"github.com/Chu16537/module_master/proto"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// 推送
func (h *Handler) Pub(subject string, data []byte) error {
	err := h.nc.Publish(subject, data)
	if err != nil {
		// 錯誤是 沒有Subject
		if err.Error() != "nats: no response from stream" {
			return err
		}

		//  創建subject
		err = h.createSubjects(subject)
		if err != nil {
			return err
		}

		err = h.nc.Publish(subject, data)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

// 訂閱
func (h *Handler) Sub(subjectName string, mode SubMode, subChan chan proto.MQSubData) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	err := h.createSubjects(subjectName)
	if err != nil {
		return err
	}

	opts := []nats.SubOpt{}

	switch mode.Mode {
	case Sub_Mode_Last_Ack:
		opts = append(opts, nats.Durable(subjectName)) // 從最後ack 開始
	case Sub_Mode_Last:
		opts = append(opts, nats.DeliverNew()) // 從訂閱後的最新消息開始
	case Sub_Mode_Sequence:
		if mode.StartSequenceID < 1 {
			mode.StartSequenceID = 1
		}
		opts = append(opts, nats.StartSequence(mode.StartSequenceID)) // 從指定的 Sequence 開始
	}

	// 使用推送型訂閱
	sub, err := h.js.Subscribe(subjectName, func(msg *nats.Msg) {
		d := proto.MQSubData{
			Data: msg.Data,
		}

		data, err := msg.Metadata()
		if err != nil {
			d.SequenceID = 0
		} else {
			d.SequenceID = data.Sequence.Consumer
		}

		subChan <- d
		msg.Ack()
	}, opts...)

	if err != nil {
		return errors.New(err.Error())
	}

	// 記錄訂閱
	h.subMap[subjectName] = sub

	return nil
}

// 取消訂閱
func (h *Handler) UnSub(subjectName string) {
	err := h.subMap[subjectName].Unsubscribe()
	if err != nil {
		fmt.Println("UnSub err", err)
	}
	delete(h.subMap, subjectName)
}
