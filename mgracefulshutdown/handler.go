package mgracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	WaitTime int // 等待時間(秒)
}

type handler struct {
	ctx          context.Context
	cancel       context.CancelFunc
	shutdownChan chan os.Signal
	mu           sync.RWMutex
	FuncMap      map[int][]func()
	MaxWaitTime  time.Duration
	ctxMap       map[int]context.Context
	cancelMap    map[int]context.CancelFunc
}

var h *handler

// 初始化
/*
waitTime 等待時間
*/
func Init(conf *Config) {
	h = &handler{
		FuncMap:     make(map[int][]func()),
		MaxWaitTime: time.Duration(conf.WaitTime) * time.Second,
	}

	h.ctx, h.cancel = context.WithCancel(context.Background())
	h.ctxMap = make(map[int]context.Context)
	h.cancelMap = make(map[int]context.CancelFunc)

	h.shutdownChan = make(chan os.Signal, 1)
	signal.Notify(h.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go shutdown()
}

// 關閉流程
func shutdown() {
	<-h.shutdownChan

	// 创建一个新的 context，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), h.MaxWaitTime)
	defer cancel()

	go execute()

	// 等待 h.MaxWaitTime 會往下執行
	<-ctx.Done()
	h.cancel()
}

// 執行關閉func
func execute() {
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
	if h == nil {
		panic("mgracefulshutdown nil")
	}

	<-h.ctx.Done()
}

// 新增 shutdown 要執行的func
func AddshutdownFunc(level int, f func()) {
	if h == nil {
		panic("mgracefulshutdown nil")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.FuncMap[level]; !ok {
		h.FuncMap[level] = make([]func(), 1)
		h.FuncMap[level][0] = f
	} else {
		h.FuncMap[level] = append(h.FuncMap[level], f)
	}

	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}
}

// 取得指定的cxt
func GetLevelCxt(level int) (context.Context, context.CancelFunc) {
	if h == nil {
		panic("mgracefulshutdown nil")
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}

	return h.ctxMap[level], h.cancelMap[level]
}

// 退出
func Exit(id int) {
	os.Exit(id)
}
