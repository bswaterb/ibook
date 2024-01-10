package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"ibook/pkg/middlewares/jwtauth"
	"ibook/pkg/utils/request"
	"ibook/pkg/utils/result"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	svc         service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
}

func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^[A-Za-z0-9]+([_\\.][A-Za-z0-9]+)*@([A-Za-z0-9\\-]+\\.)+[A-Za-z]{2,6}$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		phoneRegexPattern    = `^1[3456789]\d{9}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
	}
}

// RegisterRoutesV1 !#登录接口字段需对齐短信/手机号登录
func (u *UserHandler) RegisterRoutesV1(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile/:userId", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) SignUp(context *gin.Context) {
	req := &UserSignupReq{}
	if err1, err2 := request.ParseRequestBody(context, req), ValidateUserSignupReq(req); err1 != nil && err2 != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}

	if ok, _ := u.emailExp.MatchString(req.Email); !ok {
		result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "邮箱或密码格式错误", nil)
		return
	}
	if ok, _ := u.passwordExp.MatchString(req.Password); !ok {
		result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "邮箱或密码格式错误", nil)
		return
	}
	user, err := u.svc.SignUp(context, req.Email, req.NickName, req.Password, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.UserAlreadyExistsErr) {
			result.RespWithError(context, result.RECORD_ALREADY_EXISTS_CODE, "此邮箱已被注册", nil)
			return
		} else if errors.Is(err, service.PasswordNotEqualErr) {
			result.RespWithError(context, result.TWO_PASSWORD_NOT_EQUAL_CODE, "两次输入的密码不一致", nil)
			return
		}
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "服务内部异常，请联系管理员", nil)
		log.Println(err)
		return
	}
	result.RespWithSuccess(context, "注册成功", &UserSignupResp{
		UserId:   user.Id,
		Email:    user.Email,
		NickName: user.NickName,
		Token:    jwtauth.GenerateToken(user.Id),
	})
}

func (u *UserHandler) Login(context *gin.Context) {
	req := &UserLoginReq{}
	if err1, err2 := request.ParseRequestBody(context, req), ValidateUserLoginReq(req); err1 != nil && err2 != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if email := strings.TrimSpace(req.Email); email != "" {
		req.Email = email
		if ok, _ := u.emailExp.MatchString(req.Email); !ok {
			result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "邮箱或密码格式错误", nil)
			return
		}
	} else if phoneNumber := strings.TrimSpace(req.PhoneNumber); phoneNumber != "" {
		req.PhoneNumber = phoneNumber
		if ok, _ := u.phoneExp.MatchString(req.PhoneNumber); !ok {
			result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "手机号码或密码格式错误", nil)
			return
		}
	} else {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if ok, _ := u.passwordExp.MatchString(req.Password); !ok {
		result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "账号或密码格式错误", nil)
		return
	}
	user, err := u.svc.Login(context, req.PhoneNumber, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.UserNotExistsErr) || errors.Is(err, service.PasswordNotRightErr) {
			result.RespWithError(context, result.EMAIL_OR_PASSWORD_ERROR_CODE, "用户名或密码错误", nil)
			return
		}
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "服务内部异常，请联系管理员", nil)
	}
	result.RespWithSuccess(context, "登录成功", &UserLoginResp{
		UserId:   user.Id,
		NickName: user.NickName,
		Email:    user.Email,
		Token:    jwtauth.GenerateToken(user.Id),
	})
}

func (u *UserHandler) Profile(context *gin.Context) {
	userId := context.Param("userId")
	if strings.TrimSpace(userId) == "" {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	uId, err := strconv.Atoi(userId)
	if err != nil || uId <= 0 {
		result.RespWithError(context, result.PARAM_NOT_FULL_CODE, "请求传参或设置有误", nil)
		return
	}

	res, err := u.svc.Profile(context, int64(uId))
	if err != nil {
		if errors.Is(err, service.UserNotExistsErr) {
			result.RespWithError(context, result.USER_DO_NOT_EXISTS_CODE, "用户不存在", nil)
			return
		}
		log.Println(err)
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "服务内部异常，请联系管理员", nil)
		return
	}
	result.RespWithSuccess(context, "获取成功", &UserProfileResp{
		UserId:   res.Id,
		Email:    res.Email,
		NickName: res.NickName,
	})
}

func (u *UserHandler) Edit(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusAccepted,
		"msg":  "功能待完善",
	})
}

func (u *UserHandler) SendLoginSMSCode(context *gin.Context) {
	req := &UserSmsLoginSendReq{}
	if err1, err2 := request.ParseRequestBody(context, req), ValidateUserSmsLoginSendReq(req); err1 != nil && err2 != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if ok, _ := u.phoneExp.MatchString(req.PhoneNumber); !ok {
		result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "手机号码格式输入错误", nil)
		return
	}
	err := u.svc.SendLoginVerifyCode(context, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, service.VerifyCodeSendTooManyErr) {
			result.RespWithError(context, result.VERIFY_CODE_SEND_TOO_MANY_CODE, "验证码发送过于频繁", nil)
		} else {
			result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "系统错误", nil)
		}
		return
	}
	result.RespWithSuccess(context, "发送成功", nil)
}

func (u *UserHandler) LoginSMS(context *gin.Context) {
	req := &UserSmsLoginReq{}
	if err1, err2 := request.ParseRequestBody(context, req), ValidateUserSmsLoginReq(req); err1 != nil && err2 != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if ok, _ := u.phoneExp.MatchString(req.PhoneNumber); !ok {
		result.RespWithError(context, result.PARAM_FORMAT_ERROR_CODE, "手机号码格式有误", nil)
		return
	}
	res, err := u.svc.LoginSMS(context, req.PhoneNumber, req.VerifyCode)
	if err != nil {
		// 1. 验证码校验环节 - 1.1 验证码存在，但用户输入错误  1.2 手机号当前不存在有效验证码 1.2 校验服务内部错误
		if errors.Is(err, service.VerifyCodeComparedErr) {
			result.RespWithError(context, result.VERIFY_CODE_COMPARED_ERROR_CODE, "验证码输入错误，请尝试再次输入", nil)
		} else if errors.Is(err, service.VerifyCodeNotExists) {
			result.RespWithError(context, result.VERIFY_CODE_NOT_EXISTS_CODE, "当前手机号不存在有效验证码", nil)
		} else if errors.Is(err, service.VerifyCodeRetryTooManyErr) {
			result.RespWithError(context, result.VERIFY_CODE_RETRY_TOO_MANY_CODE, "重试次数过多，请尝试获取新的验证码", nil)
		} else if errors.Is(err, service.VerifyCodeUnKnownErr) {
			result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "未知错误，请联系管理员", nil)
			// 接下来判断错误是否是用户服务部分的问题 2.1 用户存在，查找过程出错  2.2 用户不存在，创建过程出错 2.3 用户已存在，但仍试图创建
		} else if errors.Is(err, service.UserAlreadyExistsErr) {
			result.RespWithError(context, result.RECORD_ALREADY_EXISTS_CODE, "该手机号已绑定用户，无法重复创建", nil)
		} else if errors.Is(err, service.UserUnKnownErr) || errors.Is(err, service.UserCreateErr) {
			result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "未知错误，请联系管理员", nil)
		} else {
			result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "未知错误，请联系管理员", nil)
		}
		return
	}
	result.RespWithSuccess(context, "登录成功", &UserLoginResp{
		UserId:      res.Id,
		Email:       res.Email,
		Token:       jwtauth.GenerateToken(res.Id),
		PhoneNumber: res.PhoneNumber,
	})
}
