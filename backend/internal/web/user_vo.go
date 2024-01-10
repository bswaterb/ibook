package web

import "strings"

type UserSignupReq struct {
	Email           string `json:"email" binding:"required"`
	NickName        string `json:"nickName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserSignupResp struct {
	UserId   int64  `json:"userId"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserLoginReq struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type UserLoginResp struct {
	UserId      int64  `json:"userId"`
	NickName    string `json:"nickName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Token       string `json:"token"`
}

type UserProfileResp struct {
	UserId   int64  `json:"userId"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
}

type UserSmsLoginSendReq struct {
	PhoneNumber string `json:"phoneNumber"`
}

type UserSmsLoginReq struct {
	PhoneNumber string `json:"phoneNumber"`
	VerifyCode  string `json:"verifyCode"`
}

func ValidateUserSignupReq(req *UserSignupReq) error {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.ConfirmPassword) == "" {
		return InvalidReqBodyErr
	}
	return nil
}

func ValidateUserLoginReq(req *UserLoginReq) error {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return InvalidReqBodyErr
	}
	return nil
}

func ValidateUserSmsLoginSendReq(req *UserSmsLoginSendReq) error {
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return InvalidReqBodyErr
	}
	return nil
}

func ValidateUserSmsLoginReq(req *UserSmsLoginReq) error {
	if strings.TrimSpace(req.VerifyCode) == "" || strings.TrimSpace(req.PhoneNumber) == "" {
		return InvalidReqBodyErr
	}
	return nil
}
