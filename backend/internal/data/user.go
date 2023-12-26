package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"ibook/internal/service"
	"time"
)

type User struct {
	Id          int64          `gorm:"primaryKey,autoIncrement"`
	Email       sql.NullString `gorm:"unique"`
	PhoneNumber sql.NullString `gorm:"unique"`
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
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		PhoneNumber: sql.NullString{
			String: user.PhoneNumber,
			Valid:  user.PhoneNumber != "",
		},
		Password:    user.PassWord,
		CreatedTime: now,
		UpdatedTime: now,
	}
	res := ur.db.mdb.Create(u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return service.UserAlreadyExistsErr
		}
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
		Id:          u.Id,
		Email:       u.Email.String,
		PhoneNumber: u.PhoneNumber.String,
		PassWord:    u.Password,
	}, nil
}

func (ur *userRepo) FindUserById(userId int64) (*service.User, error) {
	u := &User{}
	res := ur.db.mdb.Where("id=?", userId).First(u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, service.UserNotExistsErr
		}
		return nil, res.Error
	}
	return &service.User{
		Id:          u.Id,
		Email:       u.Email.String,
		PhoneNumber: u.PhoneNumber.String,
		PassWord:    u.Password,
	}, nil
}

func (ur *userRepo) FindUserByPhone(number string) (*service.User, error) {
	u := &User{}
	res := ur.db.mdb.Where("phone_number=?", number).First(u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, service.UserNotExistsErr
		}
		return nil, res.Error
	}
	return &service.User{
		Id:          u.Id,
		Email:       u.Email.String,
		PhoneNumber: u.PhoneNumber.String,
		PassWord:    u.Password,
	}, nil
}

type userCache struct {
	rdb redis.Cmdable
}

func NewUserCache(rdb redis.Cmdable) service.UserCache {
	return &userCache{rdb: rdb}
}

func (u *userCache) GetUserById(ctx *gin.Context, userId int64) (*service.User, error) {
	key := fmt.Sprintf("user:info:%d", userId)
	user := &service.User{}
	res, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, service.UserNotInCacheErr
		default:
			return nil, err
		}
	}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u *userCache) SetUserById(ctx *gin.Context, user *service.User) error {
	key := fmt.Sprintf("user:info:%d", user.Id)
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = u.rdb.Set(ctx, key, jsonData, time.Minute*30).Result()
	if err != nil {
		return err
	}
	return nil
}
