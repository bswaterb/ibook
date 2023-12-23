package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespWithError(ctx *gin.Context, errCode int64, msg string, result any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":   errCode,
		"msg":    msg,
		"result": result,
	})
}

func RespWithSuccess(ctx *gin.Context, msg string, result any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"msg":    msg,
		"result": result,
	})
}
