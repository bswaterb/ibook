package ratelimit

import "context"

type Limiter interface {
	// Limit key 就是限流对象
	// bool 代表是否限流，true 就是要限流
	Limit(ctx context.Context, key string) (bool, error)
}
