package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"ibook/pkg/middlewares/jwtauth"
	"ibook/pkg/utils/request"
	"ibook/pkg/utils/result"
	"log"
	"net/http"
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

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
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
		result.RespWithError(context, http.StatusInternalServerError, "服务故障，注册失败", nil)
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
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusAccepted,
		"msg":  "功能待完善",
	})
}

func (u *UserHandler) Profile(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusAccepted,
		"msg":  "功能待完善",
	})
}

func (u *UserHandler) Edit(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusAccepted,
		"msg":  "功能待完善",
	})
}
