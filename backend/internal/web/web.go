package web

import (
	"errors"
	"github.com/google/wire"
)

// WebProviderSet is data providers.
var WebProviderSet = wire.NewSet(NewUserHandler)

var (
	InvalidReqBodyErr = errors.New("请求参数不合法")
)

type UserSignupReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserSignupResp struct {
	UserId int64  `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func ValidateUserSignupReq(req *UserSignupReq) error {
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return InvalidReqBodyErr
	}
	return nil
}
