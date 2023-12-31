package pkg

import (
	"github.com/google/wire"
	"ibook/pkg/utils/logger"
	"ibook/pkg/utils/ratelimit"
)

var PkgProviderSet = wire.NewSet(ratelimit.NewRedisSlidingWindowLimiter, logger.NewZapLogger)
