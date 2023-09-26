package zgracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

const defaultFuncCount = 5 // 每个级别默认的函数数量

var handler *Handler

type Handler struct {
	ctx          context.Context
	cancel       context.CancelFunc
	shutdownChan chan os.Signal
	mu           sync.Mutex
	FuncMap      map[int][]func()
	MaxWaitTime  time.Duration
}

// 初始化
/*
waitTime 等待時間
*/
func Init(waitTime int) {
	handler = &Handler{
		FuncMap:     make(map[int][]func()),
		MaxWaitTime: time.Duration(waitTime) * time.Second,
	}

	handler.ctx, handler.cancel = context.WithCancel(context.Background())

	handler.shutdownChan = make(chan os.Signal, 1)
	signal.Notify(handler.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go handler.shutdown()

}

// 關閉流程
func (h *Handler) shutdown() {
	<-handler.shutdownChan

	// 创建一个新的 context，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), h.MaxWaitTime)
	defer cancel()

	go h.execute()

	<-ctx.Done()
	handler.cancel()
}

// 執行關閉func
func (h *Handler) execute() {
	h.mu.Lock()
	defer h.mu.Unlock()

	levels := make([]int, len(h.FuncMap))
	idx := 0
	for i := range h.FuncMap {
		levels[idx] = i
		idx++
	}
	sort.Ints(levels)

	// 根據 level 執行func
	for _, level := range levels {
		for _, f := range h.FuncMap[level] {
			f()
		}
	}
}

// 等待關閉
func WaitDone() {
	<-handler.ctx.Done()
}

// 新增 shutdown 要執行的func
func AddshutdownFunc(level int, f func()) {
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if _, ok := handler.FuncMap[level]; !ok {
		handler.FuncMap[level] = make([]func(), defaultFuncCount)
	} else {
		handler.FuncMap[level] = append(handler.FuncMap[level], f)
	}
}
