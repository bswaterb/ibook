package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	SMSSendTooManyErr = errors.New("短信发送太频繁")
	SMSUnKnownErr     = errors.New("短信服务未知错误")
)

type SMSRepo interface {
	SendMessage(ctx *gin.Context, tplId string, phoneNumbers []string, args []MsgArgs) error
}

type SMSService struct {
	smsRepo SMSRepo
}

func NewSMSService(smsRepo SMSRepo) *SMSService {
	return &SMSService{smsRepo: smsRepo}
}

type MsgArgs struct {
	Name  string
	Value string
}
