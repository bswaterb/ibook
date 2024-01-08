package web

import (
	"errors"
	"github.com/google/wire"
)

// WebProviderSet is data providers.
var WebProviderSet = wire.NewSet(NewUserHandler, NewArticleHandler)

var (
	InvalidReqBodyErr = errors.New("请求参数不合法")
)
