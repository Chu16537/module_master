package mmq

const (
	Sub_Mode_Last_Ack int = iota // 從最後ack 開始 consumer 不可以重複
	Sub_Mode_Last                // 從訂閱後的最新消息開始
	Sub_Mode_Sequence            // 從指定的 Sequence 開始
)

// 訂閱模式
type SubMode struct {
	Mode            int // 0:從最後ack 開始 1:從訂閱後的最新消息開始 2:從指定的 Sequence 開始
	StartSequenceID uint64
}

type MQSubData struct {
	SequenceID uint64 `json:"sequence_id"`
	Data       []byte `json:"data"`
	Ack        func() `json:"ack"`
}

func CreateSubChan(len int) chan MQSubData {
	return make(chan MQSubData, len)
}
