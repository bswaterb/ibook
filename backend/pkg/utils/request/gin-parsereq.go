package request

import (
	"github.com/gin-gonic/gin/binding"
)
import "github.com/gin-gonic/gin"

func ParseRequestBody(ctx *gin.Context, obj any) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())
	if err := ctx.ShouldBindWith(obj, b); err != nil {
		return err
	}
	return nil
}
