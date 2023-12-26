package service

import "github.com/gin-gonic/gin"

type ShortMessageRepo interface {
	SendMessage(ctx *gin.Context, phoneNumber string, args []MsgArgs, values []MsgArgs) error
}

type MsgArgs struct {
	name  string
	value string
}
