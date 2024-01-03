package failover

import (
	"context"
	"github.com/gin-gonic/gin"
	"ibook/internal/service/message/sms"
	"sync/atomic"
)

type timeOutSMSRepo struct {
	repos     []sms.SMSRepo
	idx       int32
	cnt       int32
	threshold int32
}

func NewTimeOutSMSRepo(repos []sms.SMSRepo) sms.SMSRepo {
	return &timeOutSMSRepo{repos: repos}
}

func (t *timeOutSMSRepo) SendMessage(ctx *gin.Context, tplId string, phoneNumbers []string, args []sms.MsgArgs) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		// 这里要切换，新的下标，往后挪了一个
		newIdx := (idx + 1) % int32(len(t.repos))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 我成功往后挪了一位
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}

	repo := t.repos[idx]
	err := repo.SendMessage(ctx, tplId, phoneNumbers, args)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		// 连续状态被打断
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		// 奇奇怪怪的异常
		return err
	}
}
