package mgracefulshutdown

import (
	"sync"
)

type wg struct {
	mu        sync.RWMutex
	taskCount int
}

func NewWg(level int) *wg {
	return &wg{}
}

func (w *wg) Add() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.taskCount++
}

func (w *wg) Done() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.taskCount--
}

func (w *wg) IsDone() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.taskCount == 0
}
