package mgin

import (
	"net/http"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/gin-gonic/gin"
)

// 請求返回 timeout
// 使用 StatusOK 原因連線成功，但是業務邏輯失敗
func ResponseTimeout(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": errorcode.Timeout,
	})
}

// 請求返回 req Unmarshal 失敗
// 使用 StatusOK 原因連線成功，但是業務邏輯失敗
func ResponseUnmarshalFail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": errorcode.Data_Unmarshal_Error,
	})
}

// 請求返回
// 使用 StatusOK 原因連線成功，但是業務邏輯失敗
func Response(c *gin.Context, err *errorcode.Error, data interface{}) {
	res := gin.H{
		"code": err.Code(),
	}

	if data != nil {
		res["data"] = data
	}

	c.JSON(http.StatusOK, res)
}
