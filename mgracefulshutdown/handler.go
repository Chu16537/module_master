package mgracefulshutdown

import (
	"context"
	"fmt"
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
	MaxWaitTime  time.Duration
	funcMap      map[int][]func()
	ctxMap       map[int]context.Context
	cancelMap    map[int]context.CancelFunc
	wgMap        map[int]*wg // 每個 level 的 WaitGroup
}

var h *handler

// 初始化
/*
waitTime 等待時間
*/
func Init(conf *Config) {
	h = &handler{
		MaxWaitTime: time.Duration(conf.WaitTime) * time.Second,
		funcMap:     make(map[int][]func()),
		ctxMap:      make(map[int]context.Context),
		cancelMap:   make(map[int]context.CancelFunc),
		wgMap:       make(map[int]*wg),
	}

	h.ctx, h.cancel = context.WithCancel(context.Background())

	h.shutdownChan = make(chan os.Signal, 1)
	signal.Notify(h.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go shutdown()
}

// 關閉流程
func shutdown() {
	<-h.shutdownChan

	// 创建一个新的 context，设置超时时间
	ctx, cancel := context.WithTimeout(h.ctx, h.MaxWaitTime)
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

	levels := make([]int, len(h.funcMap))
	idx := 0
	for i := range h.funcMap {
		levels[idx] = i
		idx++
	}
	sort.Ints(levels)

	// 根據 level 執行func
	fmt.Println("開始執行關閉函數...")
	for _, level := range levels {
		if wg, exists := h.wgMap[level]; exists {
			for {
				if wg.IsDone() {
					break // 如果該 level 的任務已完成，進入下一個 level
				}
				time.Sleep(2 * time.Second) // 等待 2 秒後重新檢查
			}
		}

		fmt.Printf("執行 level: %d 的關閉函數，共 %d 個\n", level, len(h.funcMap[level]))
		for _, f := range h.funcMap[level] {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("關閉函數失敗，level: %d, error: %v\n", level, r)
					}
				}()
				f()
			}()
		}
	}
	fmt.Println("所有關閉函數執行完成。")
}

// 主動關閉
func Shutdown() {
	h.shutdownChan <- syscall.SIGTERM
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

	if _, ok := h.funcMap[level]; !ok {
		h.funcMap[level] = []func(){}
	}

	h.funcMap[level] = append(h.funcMap[level], f)

	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}

}

// 取得指定的cxt
func GetLevelCxt(level int) (context.Context, context.CancelFunc) {
	if h == nil {
		panic("mgracefulshutdown nil")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctxMap[level] == nil {
		h.ctxMap[level], h.cancelMap[level] = context.WithCancel(context.Background())
	}

	if _, ok := h.wgMap[level]; !ok {
		h.wgMap[level] = NewWg(level)
	}

	return h.ctxMap[level], h.cancelMap[level]
}

// 添加需要等待完成的任務
func AddTask(level int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.wgMap[level]; !exists {
		h.wgMap[level] = NewWg(level)
	}
	h.wgMap[level].Add()
}

// 標記任務已完成
func DoneTask(level int) {
	wg, exists := h.wgMap[level]
	if exists {
		wg.Done()
	}
}

// 退出
func Exit(id int) {
	os.Exit(id)
}
