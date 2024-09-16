package mgin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	IsSwag        bool
	Addr          string
	TimeoutSecond int
}

type Handler struct {
	ctx             context.Context
	config          *Config
	routine         *gin.Engine
	timeoutDuration time.Duration
	srv             *http.Server
	timeoutFunc     func(c *gin.Context)
}

func New(ctx context.Context, config *Config, timeoutFn func(c *gin.Context)) (*Handler, error) {
	// 基本判斷
	if config.Addr == "" {
		return nil, errors.New("gin new error addr nil")
	}

	if config.TimeoutSecond == 0 {
		config.TimeoutSecond = 5
	}

	tFn := timeoutBase
	if timeoutFn != nil {
		tFn = timeoutFn
	}

	h := &Handler{
		ctx:             ctx,
		config:          config,
		timeoutDuration: time.Duration(time.Duration(config.TimeoutSecond) * time.Second),
		timeoutFunc:     tFn,
	}

	routine := gin.Default()
	routine.Use(h.middlewareTimeout())

	// show swag
	if h.config.IsSwag {
		routine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	h.routine = routine

	return h, nil
}

func (h *Handler) Done() {
	h.srv.Shutdown(h.ctx)
}

func (h *Handler) Run() error {
	// 配置 HTTP 服务器
	h.srv = &http.Server{
		Addr:    h.config.Addr,
		Handler: h.routine,
	}

	errChan := make(chan error, 1)

	go func() {
		err := h.srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 等待n秒判斷是否有錯
	select {
	case err := <-errChan:
		return err

	case <-time.After(5 * time.Second):
		// 等待5秒發現沒有錯誤
		return nil
	}

}

func timeoutBase(c *gin.Context) {
	res := gin.H{
		"msg": "timeout",
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}
