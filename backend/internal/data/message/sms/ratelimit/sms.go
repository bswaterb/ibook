package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ibook/internal/data/message/sms/mem"
	"ibook/internal/service/message/sms"
	"ibook/pkg/utils/ratelimit"
)

type rateLimitSmsRepo struct {
	repo    sms.SMSRepo
	limiter ratelimit.Limiter
}

func NewRateLimitSmsRepo(limiter ratelimit.Limiter) sms.SMSRepo {
	return &rateLimitSmsRepo{
		repo:    mem.NewMemSMSRepo(),
		limiter: limiter,
	}
}

func (r *rateLimitSmsRepo) SendMessage(ctx *gin.Context, tplId string, phoneNumbers []string, args []sms.MsgArgs) error {
	// 限流措施
	limited, err := r.limiter.Limit(ctx, "sms:aliyun")
	if err != nil {
		return fmt.Errorf("短信限流模块出现问题，%w", err)
	}
	if limited {
		// 转异步消费
		return fmt.Errorf("短信发送服务当前处于限流状态，拒绝发送至: %s", phoneNumbers[0])
	}

	err = r.repo.SendMessage(ctx, tplId, phoneNumbers, args)
	return err
}
