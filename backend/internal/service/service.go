package service

import "github.com/google/wire"

// ServiceProviderSet is data providers.
var ServiceProviderSet = wire.NewSet(NewUserService)

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	PassWord string
}
