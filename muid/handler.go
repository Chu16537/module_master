package muid

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/rs/xid"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxSequenceVal = uint32(1<<32 - 1) // 避免sequence溢位的最大值
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
	randGen       *rand.Rand
}

var h *Handler

func New(node int64) *Handler {
	var (
		sequence     uint32  = 0
		sequenceBits uintptr = uintptr(32)
		instanceMask         = node << int64(sequenceBits)
		tailBits     uintptr = uintptr(20)
	)

	t := time.Now()

	h = &Handler{
		node:          node,
		uid:           xid.New().String(),
		time:          t,
		lastTimestamp: t.UnixNano(),
		sequence:      sequence,
		sequenceBits:  sequenceBits,
		instanceMask:  instanceMask,
		tailBits:      tailBits,
		randGen:       rand.New(rand.NewSource(t.UnixNano())),
	}

	return h
}

// 創建唯一碼
func (h *Handler) CreateID() int64 {
	nowTimestamp := time.Now().UnixNano() >> int64(h.tailBits)

	// 同一時間戳下，防止sequence溢位
	if h.lastTimestamp == nowTimestamp {
		if atomic.AddUint32(&h.sequence, 1) > maxSequenceVal {
			for nowTimestamp <= h.lastTimestamp {
				nowTimestamp = h.time.UnixNano() >> int64(h.tailBits)
			}
			h.sequence = 0
		}
	} else {
		h.sequence = 0
	}

	h.lastTimestamp = nowTimestamp

	return nowTimestamp<<int64(h.tailBits) | h.instanceMask | int64(h.sequence)
}

// 隨機字符串生成函數
func (h *Handler) CreatRandomString(length int) string {
	result := make([]byte, length)
	charsetLength := len(charset)
	for i := 0; i < length; i++ {
		result[i] = charset[h.randGen.Intn(charsetLength)]
	}
	return string(result)
}

// orderID 時間-nodeId-唯一碼
func (h *Handler) CreateOrderID() string {
	fmt.Println(h.uid)
	return fmt.Sprintf("%v-%v-%v", time.Now().Unix(), h.node, h.CreateID())
}
