package data

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"ibook/internal/service"
	"strconv"
)

//go:embed script/lua/verify_code/set_code.lua
var luaSetCode string

//go:embed script/lua/verify_code/check_code.lua
var luaCheckCode string

type verifyCodeRepo struct {
	rdb redis.Cmdable
}

func NewVerifyCodeRepo(rdb redis.Cmdable) service.VerifyCodeRepo {
	return &verifyCodeRepo{rdb: rdb}
}

func (v *verifyCodeRepo) SetVerifyCode(ctx *gin.Context, key string, verifyCode string, lifeDurationSeconds int64) error {
	res, err := v.rdb.Eval(ctx, luaSetCode, []string{key}, []string{verifyCode, strconv.FormatInt(lifeDurationSeconds, 10)}).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return service.VerifyCodeSendTooManyErr
	case -2:
		return service.VerifyCodeLifeTimeErr
	default:
		return service.VerifyCodeUnKnownErr
	}
}

func (v *verifyCodeRepo) CheckVerifyCode(ctx *gin.Context, key string, inputCode string) (bool, error) {
	res, err := v.rdb.Eval(ctx, luaCheckCode, []string{key}, inputCode).Int()
	if err != nil {
		return false, nil
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, service.VerifyCodeRetryTooManyErr
	case -2:
		return false, service.VerifyCodeComparedErr
	case -3:
		return false, service.VerifyCodeNotExists
	default:
		return false, service.VerifyCodeUnKnownErr
	}
}
