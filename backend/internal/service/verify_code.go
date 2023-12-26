package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type VerifyCodeRepo interface {
	SetVerifyCode(ctx *gin.Context, key string, verifyCode string, lifeDurationSeconds time.Duration) error
	GetVerifyCode(ctx *gin.Context, key string) (string, error)
}

type VerifyCodeService struct {
	smr ShortMessageRepo
	vcr VerifyCodeRepo
}

func (service *VerifyCodeService) SendVerifyCode(ctx *gin.Context, bizName string, phoneNumber string) error {
	// 1. 生成验证码
	code := service.generateVerifyCode()
	key := fmt.Sprintf("verify_code:%s:%s", bizName, phoneNumber)
	// 2. 将验证码存储到 Redis 中
	err := service.vcr.SetVerifyCode(ctx, key, code, time.Minute*30)
	if err != nil {
		// 有错误说明 Redis 缓存出现问题
		return err
	}
	// 3. 发送验证码短信
	err = service.smr.SendMessage(ctx, phoneNumber, []MsgArgs{}, []MsgArgs{})
	if err != nil {
		// 说明短信发送失败，但 Redis 中有缓存
		return errors.New("短信发送失败")
	}
	return nil
}

func (service *VerifyCodeService) generateVerifyCode() string {
	num := rand.Intn(1000000)
	res := fmt.Sprintf("%06d", num)
	return res
}
