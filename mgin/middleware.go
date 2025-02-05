package mgin

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddMiddleware(f func() gin.HandlerFunc) {
	h.routine.Use(f())
}

// 中间件 timeout
func (h *Handler) middlewareTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为请求创建一个带超时的 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), h.config.TimeoutSecond)
		defer cancel()

		// 使用 context 替换原始的请求 context
		c.Request = c.Request.WithContext(ctx)

		// 创建一个 channel 来监听请求是否完成
		done := make(chan struct{})
		// 使用 Goroutine 处理请求
		go func() {
			c.Next()
			defer close(done)
		}()

		select {
		case <-ctx.Done():
			h.timeoutFunc(c)
			return
		case <-done:
			// 请求正常完成，继续执行
			return
		}

	}
}
