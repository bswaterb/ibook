package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaSlideWindow string

// RedisSlidingWindowLimiter Redis 上的滑动窗口算法限流器实现
type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable

	// 窗口大小
	interval time.Duration
	// 阈值
	rate int
	// interval 内允许 rate 个请求
	// 1s 内允许 3000 个请求
}

// NewRedisSlidingWindowLimiter !#需补全配置化逻辑
func NewRedisSlidingWindowLimiter(cmd redis.Cmdable) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: time.Second,
		rate:     3000,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaSlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
