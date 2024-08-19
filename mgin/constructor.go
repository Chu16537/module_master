package mgin

import (
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

type handler struct {
	Config          *Config
	Routine         *gin.Engine
	TimeoutDuration time.Duration
}

var h *handler

func New(config *Config) error {
	// 基本判斷
	if config.Addr == "" {
		return errors.New("gin new error addr nil")
	}

	if config.TimeoutSecond == 0 {
		config.TimeoutSecond = 5
	}

	h = &handler{
		Config:          config,
		TimeoutDuration: time.Duration(time.Duration(config.TimeoutSecond) * time.Second),
	}

	routine := gin.Default()
	routine.Use(middlewareTimeout(h.TimeoutDuration))

	// show swag
	if h.Config.IsSwag {
		routine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	h.Routine = routine

	return nil
}

func Get() *handler {
	return h
}

func Run() error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- h.Routine.Run(h.Config.Addr)
	}()

	err := <-errChan

	if err != nil {
		return err
	}

	return nil
}

func Done() {

}
