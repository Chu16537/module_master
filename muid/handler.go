package muid

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxSequenceVal = uint32(1<<32 - 1) // 避免sequence溢位的最大值
)

type handler struct {
	node int64 // 機器碼

	time          time.Time // 時間
	lastTimestamp int64     // 上次使用時間
	sequence      uint32    // 順序號碼
	instanceMask  int64
	tailBits      uintptr
	randGen       *rand.Rand
	orderIDCount  uint32 // 使用次數
}

var h *handler

func New(node int64) {
	t := time.Now()

	h = &handler{
		node:          node,
		time:          t,
		lastTimestamp: t.UnixNano(),
		sequence:      0,
		instanceMask:  node << int64(uintptr(32)),
		tailBits:      uintptr(20),
		randGen:       rand.New(rand.NewSource(t.UnixNano())),
		orderIDCount:  0,
	}
}

func GetNodeID() int64 {
	return h.node
}

// 創建唯一碼
func CreateID() int64 {
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
func CreatRandomString(length int) string {
	result := make([]byte, length)
	charsetLength := len(charset)
	for i := 0; i < length; i++ {
		result[i] = charset[h.randGen.Intn(charsetLength)]
	}
	return string(result)
}

// orderID 時間-nodeId-唯一碼
func CreateOrderID() string {
	if atomic.AddUint32(&h.orderIDCount, 1) > maxSequenceVal {
		h.orderIDCount = 0
	}

	return fmt.Sprintf("%v-%v-%v", CreateID(), h.node, h.orderIDCount)
}
