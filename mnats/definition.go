package mnats

// 訂閱模式
type SubMode struct {
	Mode            int // 0:從最後ack 開始 1:從訂閱後的最新消息開始 2:從指定的 Sequence 開始
	StartSequenceID uint64
}

const (
	Sub_Mode_Last_Ack int = iota // 從最後ack 開始
	Sub_Mode_Last                // 從訂閱後的最新消息開始
	Sub_Mode_Sequence            // 從指定的 Sequence 開始
)
