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

type Igf interface {
	// 取得指定的cxt
	GetLevelCxt(level int) (context.Context, context.CancelFunc)
	// 新增 shutdown 要執行的func
	AddshutdownFunc(level int, f func())
}

type Config struct {
	WaitTime int // 等待時間
}

type Handler struct {
	ctx          context.Context
	cancel       context.CancelFunc
	shutdownChan chan os.Signal
	mu           sync.Mutex
	FuncMap      map[int][]func()
	MaxWaitTime  time.Duration
	ctxMap       map[int]context.Context
	cancelMap    map[int]context.CancelFunc
}

// 初始化
/*
waitTime 等待時間
*/
func Init(conf *Config) *Handler {
	h := &Handler{
		FuncMap:     make(map[int][]func()),
		MaxWaitTime: time.Duration(conf.WaitTime) * time.Second,
	}

	h.ctx, h.cancel = context.WithCancel(context.Background())
	h.ctxMap = make(map[int]context.Context)
	h.cancelMap = make(map[int]context.CancelFunc)

	h.shutdownChan = make(chan os.Signal, 1)
	signal.Notify(h.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go h.shutdown()

	return h
}

// 關閉流程
func (h *Handler) shutdown() {
	<-h.shutdownChan

	// 创建一个新的 context，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), h.MaxWaitTime)
	defer cancel()

	go h.execute()

	// 等待 h.MaxWaitTime 會往下執行
	<-ctx.Done()
	h.cancel()
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
func (h *Handler) WaitDone() {
	<-h.ctx.Done()
}

// 新增 shutdown 要執行的func
func (h *Handler) AddshutdownFunc(level int, f func()) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.FuncMap[level]; !ok {
		h.FuncMap[level] = make([]func(), defaultFuncCount)
	} else {
		h.FuncMap[level] = append(h.FuncMap[level], f)
	}

	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}
}

// 取得指定的cxt
func (h *Handler) GetLevelCxt(level int) (context.Context, context.CancelFunc) {
	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}

	return h.ctxMap[level], h.cancelMap[level]
}
