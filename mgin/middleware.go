package mgin

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// 中间件 timeout
func middlewareTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为请求创建一个带超时的 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 使用 context 替换原始的请求 context
		c.Request = c.Request.WithContext(ctx)

		// 创建一个 channel 来监听请求是否完成
		done := make(chan struct{})

		go func() {
			// 调用下一个中间件/处理器
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			ResponseTimeout(c)
		case <-done:
			// 请求正常完成並回覆使用者，继续执行

		}
	}
}
