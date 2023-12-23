package data

import (
	"errors"
	"gorm.io/gorm"
	"ibook/internal/service"
	"time"
)

type User struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Email       string `gorm:"unique"`
	Password    string
	CreatedTime int64
	UpdatedTime int64
}

type userRepo struct {
	db *Data
}

func NewUserRepo(db *Data) service.UserRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) CreateUser(user *service.User) error {
	now := time.Now().UTC().UnixMilli()
	u := &User{
		Email:       user.Email,
		Password:    user.PassWord,
		CreatedTime: now,
		UpdatedTime: now,
	}
	res := ur.db.mdb.Create(u)
	if res.Error != nil {

	}
	user.Id = u.Id
	return nil
}

func (ur *userRepo) FindUserByEmail(email string) (*service.User, error) {
	u := &User{}
	res := ur.db.mdb.Where("email=?", email).First(u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, service.UserNotExistsErr
		}
		return nil, res.Error
	}
	return &service.User{
		Id:    u.Id,
		Email: u.Email,
	}, nil
}
