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
)

type UserRepo interface {
	CreateUser(user *User) error
	FindUserByEmail(email string) (*User, error)
}

type UserService struct {
	ur UserRepo
}

func NewUserService(ur UserRepo) *UserService {
	return &UserService{ur: ur}
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

func (s *UserService) Login(context *gin.Context, email string, password string) (*User, error) {
	user, err := s.ur.FindUserByEmail(email)
	if err != nil && errors.Is(err, UserNotExistsErr) {
		return nil, UserNotExistsErr
	}
	if user.PassWord != password {
		return nil, PasswordNotRightErr
	}
	return user, nil
}
