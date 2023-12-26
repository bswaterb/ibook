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
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutesV1(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile/:userId", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(context *gin.Context) {
	req := &UserSignupReq{}
	if err := request.ParseRequestBody(context, req); err != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if err := ValidateUserSignupReq(req); err != nil {
		result.RespWithError(context, result.PARAM_NOT_FULL_CODE, "请求传参或设置有误", nil)
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
	user, err := u.svc.SignUp(context, req.Email, req.Password, req.ConfirmPassword)
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
		UserId: user.Id,
		Email:  user.Email,
		Token:  jwtauth.GenerateToken(user.Id),
	})
}

func (u *UserHandler) Login(context *gin.Context) {
	req := &UserLoginReq{}
	if err := request.ParseRequestBody(context, req); err != nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if err := ValidateUserLoginReq(req); err != nil {
		result.RespWithError(context, result.PARAM_NOT_FULL_CODE, "请求传参或设置有误", nil)
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
	user, err := u.svc.Login(context, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.UserNotExistsErr) || errors.Is(err, service.PasswordNotRightErr) {
			result.RespWithError(context, result.EMAIL_OR_PASSWORD_ERROR_CODE, "用户名或密码错误", nil)
			return
		}
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "服务内部异常，请联系管理员", nil)
	}
	result.RespWithSuccess(context, "登录成功", &UserLoginResp{
		UserId: user.Id,
		Email:  user.Email,
		Token:  jwtauth.GenerateToken(user.Id),
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
		UserId: res.Id,
		Email:  res.Email,
	})
}

func (u *UserHandler) Edit(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusAccepted,
		"msg":  "功能待完善",
	})
}
