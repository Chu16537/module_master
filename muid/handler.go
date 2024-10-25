package muid

import (
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/Chu16537/module_master/mtime"
	"github.com/rs/xid"
)

const (
	// 定義字母和數字的字元集
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Handler struct {
	node          int64     // 機器碼
	uid           string    // 唯一號碼
	time          time.Time // 時間
	lastTimestamp int64     // 上次使用時間
	sequence      uint32    // 順序號碼
	sequenceBits  uintptr
	instanceMask  int64
	tailBits      uintptr
}

var h *Handler

func New(node int64) *Handler {

	var (
		sequence     uint32  = 0
		sequenceBits uintptr = unsafe.Sizeof(sequence) * 8
		instanceMask         = node << int64(sequenceBits)
		tailBits     uintptr = sequenceBits
	)

	h = &Handler{
		node:          node,
		uid:           xid.New().String(),
		time:          mtime.GetZero(),
		lastTimestamp: mtime.GetZero().UnixNano(),
		sequence:      sequence,
		sequenceBits:  sequenceBits,
		instanceMask:  instanceMask,
		tailBits:      tailBits,
	}

	return h
}

// 創建唯一碼
func (h *Handler) CreateID() int64 {
	nowTimestamp := time.Now().UnixNano() >> int64(h.tailBits)

	// 同一時間創建增加sequence
	if h.lastTimestamp == nowTimestamp {
		atomic.AddUint32(&h.sequence, 1)
	} else {
		h.sequence = 0
	}

	h.lastTimestamp = nowTimestamp

	return nowTimestamp<<int64(h.tailBits) | h.instanceMask | int64(h.sequence)

}

// 隨機字符串生成函數
func (h *Handler) CreatRandomString(length int) string {
	// 預先生成隨機數的 buffer
	seededRand := rand.New(rand.NewSource(h.CreateID()))
	result := make([]byte, length)
	buffer := make([]byte, 10) // 用來批量生成隨機數字節
	charsetLength := len(charset)

	for i := 0; i < length; i += 10 {
		// 批量生成隨機數
		seededRand.Read(buffer)

		// 將隨機數轉換為對應的字母和數字
		for j := 0; j < 10 && i+j < length; j++ {
			// 這裡使用 buffer[j] % charsetLength 來確保隨機數在字符集範圍內
			result[i+j] = charset[int(buffer[j])%charsetLength]
		}
	}
	return string(result)
}
