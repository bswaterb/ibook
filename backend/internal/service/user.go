package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	PasswordNotEqualErr  = errors.New("两次输入密码不一致")
	UserAlreadyExistsErr = errors.New("该用户已存在")
	UserNotExistsErr     = errors.New("该用户不存在")
	UserCreateErr        = errors.New("用户创建失败，请联系管理员")
	PasswordNotRightErr  = errors.New("密码不符")
	UserNotInCacheErr    = errors.New("用户信息未被缓存")
)

type UserRepo interface {
	CreateUser(user *User) error
	FindUserByEmail(email string) (*User, error)
	FindUserById(userId int64) (*User, error)
}

type UserCache interface {
	GetUserById(ctx *gin.Context, userId int64) (*User, error)
	SetUserById(ctx *gin.Context, user *User) error
}

type UserService struct {
	ur UserRepo
	uc UserCache
}

func NewUserService(ur UserRepo, uc UserCache) *UserService {
	return &UserService{ur: ur, uc: uc}
}

func (s *UserService) SignUp(ctx *gin.Context, email string, password string, confirmPassword string) (*User, error) {
	if password != confirmPassword {
		return nil, PasswordNotEqualErr
	}
	_, err := s.ur.FindUserByEmail(email)
	if err == nil || !errors.Is(err, UserNotExistsErr) {
		return nil, UserAlreadyExistsErr
	}
	user := &User{
		Email:    email,
		PassWord: password,
	}
	if err := s.ur.CreateUser(user); err != nil {
		if errors.Is(err, UserAlreadyExistsErr) {
			return nil, UserAlreadyExistsErr
		}
		return nil, UserCreateErr
	}
	return user, nil
}

func (s *UserService) Login(ctx *gin.Context, email string, password string) (*User, error) {
	user, err := s.ur.FindUserByEmail(email)
	if err != nil && errors.Is(err, UserNotExistsErr) {
		return nil, UserNotExistsErr
	}
	if user.PassWord != password {
		return nil, PasswordNotRightErr
	}
	return user, nil
}

func (s *UserService) Profile(ctx *gin.Context, userId int64) (*User, error) {
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
