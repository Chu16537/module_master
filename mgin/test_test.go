package mgin_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mgin"
	"github.com/gin-gonic/gin"
)

func TestXxx(t *testing.T) {

	config := &mgin.Config{
		Port:    "8080",
		Timeout: 10 * time.Second,
	}

	h, err := mgin.New(context.Background(), config, nil)

	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}

	r := h.GetRoutine()

	r.POST("/a", func(c *gin.Context) {
		fmt.Println("a")
	})

	err = h.Run()
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}

	time.Sleep(10 * time.Second)
}
