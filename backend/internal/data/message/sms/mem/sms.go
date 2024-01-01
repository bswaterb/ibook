package mem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ibook/internal/service/message/sms"
)

// 实现短信发送接口的 Demo
type memSMSRepo struct {
}

func NewMemSMSRepo() sms.SMSRepo {
	return &memSMSRepo{}
}

func (s *memSMSRepo) SendMessage(ctx *gin.Context, tplId string, phoneNumber []string, args []sms.MsgArgs) error {
	fmt.Println("[模拟短信发送] -> " + phoneNumber[0] + ": " + args[0].Value)
	return nil
}
