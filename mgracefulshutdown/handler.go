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

var GF *Handler

// 初始化
/*
waitTime 等待時間
*/
func Init(conf *Config) {
	GF = &Handler{
		FuncMap:     make(map[int][]func()),
		MaxWaitTime: time.Duration(conf.WaitTime) * time.Second,
	}

	GF.ctx, GF.cancel = context.WithCancel(context.Background())
	GF.ctxMap = make(map[int]context.Context)
	GF.cancelMap = make(map[int]context.CancelFunc)

	GF.shutdownChan = make(chan os.Signal, 1)
	signal.Notify(GF.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go shutdown()
}

// 關閉流程
func shutdown() {
	<-GF.shutdownChan

	// 创建一个新的 context，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), GF.MaxWaitTime)
	defer cancel()

	go execute()

	// 等待 h.MaxWaitTime 會往下執行
	<-ctx.Done()
	GF.cancel()
}

// 執行關閉func
func execute() {
	GF.mu.Lock()
	defer GF.mu.Unlock()

	levels := make([]int, len(GF.FuncMap))
	idx := 0
	for i := range GF.FuncMap {
		levels[idx] = i
		idx++
	}
	sort.Ints(levels)

	// 根據 level 執行func
	for _, level := range levels {
		for _, f := range GF.FuncMap[level] {
			f()
		}
	}
}

// 等待關閉
func WaitDone() {
	<-GF.ctx.Done()
}

// 新增 shutdown 要執行的func
func AddshutdownFunc(level int, f func()) {
	GF.mu.Lock()
	defer GF.mu.Unlock()

	if _, ok := GF.FuncMap[level]; !ok {
		GF.FuncMap[level] = make([]func(), 1)
		GF.FuncMap[level][0] = f
	} else {
		GF.FuncMap[level] = append(GF.FuncMap[level], f)
	}

	if GF.ctxMap[level] == nil {
		GF.ctxMap[level], GF.cancelMap[level] = context.WithCancel(context.Background())
	}
}

// 取得指定的cxt
func GetLevelCxt(level int) (context.Context, context.CancelFunc) {
	if GF.ctxMap[level] == nil {
		GF.ctxMap[level], GF.cancelMap[level] = context.WithCancel(context.Background())
	}

	return GF.ctxMap[level], GF.cancelMap[level]
}

// 退出
func Exit(id int) {
	os.Exit(id)
}
