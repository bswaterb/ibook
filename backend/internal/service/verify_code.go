package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	VerifyCodeLifeTimeErr     = errors.New("验证码过期时间不符合预期")
	VerifyCodeSendTooManyErr  = errors.New("验证码发送过于频繁")
	VerifyCodeRetryTooManyErr = errors.New("验证码尝试校验次数过多")
	VerifyCodeComparedErr     = errors.New("验证码输入错误")
	VerifyCodeNotExists       = errors.New("验证码不存在")
	VerifyCodeUnKnownErr      = errors.New("验证码服务未知错误")
)

type VerifyCodeRepo interface {
	SetVerifyCode(ctx *gin.Context, key string, verifyCode string, lifeDurationSeconds int64) error
	CheckVerifyCode(ctx *gin.Context, key string, inputCode string) (bool, error)
}

//
//type VerifyCodeService struct {
//	smr SMSRepo
//	vcr VerifyCodeRepo
//}
//
//func (service *VerifyCodeService) SendVerifyCode(ctx *gin.Context, bizName string, phoneNumber string) error {
//	// 1. 生成验证码
//	code := service.generateVerifyCode()
//	key := fmt.Sprintf("verify_code:%s:%s", bizName, phoneNumber)
//	// 2. 将验证码存储到 Redis 中
//	err := service.vcr.SetVerifyCode(ctx, key, code, 60*30)
//	if err != nil {
//		// 有错误说明 Redis 缓存出现问题
//		return err
//	}
//	// 3. 发送验证码短信
//	err = service.smr.SendMessage(ctx, "", []string{phoneNumber}, []MsgArgs{})
//	if err != nil {
//		// 说明短信发送失败，但 Redis 中有缓存
//		return errors.New("短信发送失败")
//	}
//	return nil
//}
//
//func (service *VerifyCodeService) generateVerifyCode() string {
//	num := rand.Intn(1000000)
//	res := fmt.Sprintf("%06d", num)
//	return res
//}
