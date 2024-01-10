package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"ibook/internal/service/message/sms"
	"ibook/pkg/utils/randcode"
)

var (
	PasswordNotEqualErr  = errors.New("两次输入密码不一致")
	UserAlreadyExistsErr = errors.New("该用户已存在")
	UserNotExistsErr     = errors.New("该用户不存在")
	UserCreateErr        = errors.New("用户创建失败，请联系管理员")
	PasswordNotRightErr  = errors.New("密码不符")
	UserNotInCacheErr    = errors.New("用户信息未被缓存")
	UserUnKnownErr       = errors.New("用户服务未知错误")
)

type UserRepo interface {
	CreateUser(user *User) error
	FindUserByEmail(email string) (*User, error)
	FindUserById(userId int64) (*User, error)
	FindUserByPhone(number string) (*User, error)
}

type UserCache interface {
	GetUserById(ctx *gin.Context, userId int64) (*User, error)
	SetUserById(ctx *gin.Context, user *User) error
}

type UserService interface {
	SignUp(ctx *gin.Context, email string, nickName string, password string, confirmPassword string) (*User, error)
	Login(ctx *gin.Context, phone string, email string, password string) (*User, error)
	Profile(ctx *gin.Context, userId int64) (*User, error)
	SendLoginVerifyCode(ctx *gin.Context, phoneNumber string) error
	LoginSMS(context *gin.Context, phoneNumber string, code string) (*User, error)
}

type userService struct {
	ur   UserRepo
	smsr sms.SMSRepo
	vcr  VerifyCodeRepo
	uc   UserCache
}

func NewUserService(ur UserRepo, uc UserCache, smsr sms.SMSRepo, vcr VerifyCodeRepo) UserService {
	return &userService{ur: ur, uc: uc, smsr: smsr, vcr: vcr}
}

func (s *userService) SignUp(ctx *gin.Context, email string, nickName string, password string, confirmPassword string) (*User, error) {
	if password != confirmPassword {
		return nil, PasswordNotEqualErr
	}
	_, err := s.ur.FindUserByEmail(email)
	if err == nil || !errors.Is(err, UserNotExistsErr) {
		return nil, UserAlreadyExistsErr
	}
	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Email:    email,
		NickName: nickName,
		PassWord: string(cryptPassword),
	}
	if err := s.ur.CreateUser(user); err != nil {
		if errors.Is(err, UserAlreadyExistsErr) {
			return nil, UserAlreadyExistsErr
		}
		return nil, UserCreateErr
	}
	return user, nil
}

func (s *userService) Login(ctx *gin.Context, phone string, email string, password string) (*User, error) {
	var user *User
	var err error
	if email != "" {
		user, err = s.ur.FindUserByEmail(email)
		if err != nil && errors.Is(err, UserNotExistsErr) {
			return nil, UserNotExistsErr
		}
	} else if phone != "" {
		user, err = s.ur.FindUserByPhone(phone)
		if err != nil && errors.Is(err, UserNotExistsErr) {
			return nil, UserNotExistsErr
		}
	}
	if user == nil {
		return nil, UserNotExistsErr
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password)) != nil {
		return nil, PasswordNotRightErr
	}
	return user, nil
}

func (s *userService) Profile(ctx *gin.Context, userId int64) (*User, error) {
	userCache, err := s.uc.GetUserById(ctx, userId)
	if err == nil {
		userCache.PassWord = ""
		return userCache, nil
	}
	if !errors.Is(err, UserNotInCacheErr) {
		return nil, err
	}
	// cache 中不存在用户的缓存记录，去 repo 中查找
	user, err := s.ur.FindUserById(userId)
	if err != nil {
		if errors.Is(err, UserNotExistsErr) {
			return nil, UserNotExistsErr
		}
		return nil, err
	}
	// 将本次查询的 user 置入缓存
	err = s.uc.SetUserById(ctx, user)
	// cache 插入失败，可考虑进一步策略
	if err != nil {
		return nil, err
	}
	user.PassWord = ""
	return user, nil
}

func (s *userService) SendLoginVerifyCode(ctx *gin.Context, phoneNumber string) error {
	code := randcode.GenVerifyCode(6, randcode.TYPE_MIXED)
	phoneNumbers := []string{phoneNumber}
	// 尝试将 code 置入缓存中
	key := fmt.Sprintf("verify_code:%s:%s", "login", phoneNumber)
	err := s.vcr.SetVerifyCode(ctx, key, code, 60*30)
	if err != nil {
		// 1. 本次与上次获取缓存操作间隔小于 1 分钟
		// 2. 缓存中过期时间设置有误
		// 3. 其他未知错误
		return err
	}
	return s.smsr.SendMessage(ctx, "verify_code_tlpId", phoneNumbers, []sms.MsgArgs{{Name: "code", Value: code}})
}

func (s *userService) LoginSMS(context *gin.Context, phoneNumber string, code string) (*User, error) {
	// 1. 校验验证码是否正确
	key := fmt.Sprintf("verify_code:%s:%s", "login", phoneNumber)
	ok, err := s.vcr.CheckVerifyCode(context, key, code)
	if err != nil {
		if !errors.Is(err, VerifyCodeComparedErr) {
			if errors.Is(err, VerifyCodeRetryTooManyErr) {
				return nil, VerifyCodeRetryTooManyErr
			} else if errors.Is(err, VerifyCodeUnKnownErr) {
				return nil, err
			}
			return nil, err
		}
	}
	// 本次输入错误
	if !ok {
		return nil, VerifyCodeComparedErr
	}
	// 2. 进入数据库获取当前登录用户的信息，若不存在则创建新用户
	u, err := s.ur.FindUserByPhone(phoneNumber)
	// 出错情况: 1. 用户不存在 2. 查找过程出错
	if err != nil {
		// 当前用户不存在，尝试创建新用户
		if errors.Is(err, UserNotExistsErr) {
			u = &User{PhoneNumber: phoneNumber}
			err = s.ur.CreateUser(u)
			// 出错情况: 1. 用户已存在 2. 创建过程出错
			if err != nil {
				// 当前用户已存在，无法重复创建 ---- 并发场景下有概率发生
				if errors.Is(err, UserAlreadyExistsErr) {
					return nil, UserAlreadyExistsErr
				} else {
					return nil, UserUnKnownErr
				}
			}
		} else {
			return nil, UserUnKnownErr
		}
	}
	return &User{
		Id:          u.Id,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		PassWord:    u.PassWord,
	}, nil
}
