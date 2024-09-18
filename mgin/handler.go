package mgin

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRoutine() *gin.Engine {
	return h.routine
}

func GetToken(c *gin.Context) string {
	return c.Request.Header.Get("X-Token")
}

func GetRequestID(c *gin.Context) string {
	return c.Request.Header.Get("X-RequestID")
}

func GetAllPath(c *gin.Context) string {
	return c.Request.URL.Path
}

// 確認最後的路徑是否正確
func IsLastPath(c *gin.Context, s string) bool {
	return strings.HasSuffix(c.Request.URL.Path, s)
}

func GetCtx(c *gin.Context) context.Context {
	return c.Request.Context()
}

func Response(c *gin.Context, errCode int, data interface{}) {
	res := gin.H{
		"code": errCode,
	}

	if data != nil {
		res["data"] = data
	}

	c.JSON(http.StatusOK, res)
}

func ResponseFail(c *gin.Context, errCode int) {
	res := gin.H{
		"code": errCode,
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}
