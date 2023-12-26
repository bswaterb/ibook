package service

import "github.com/google/wire"

// ServiceProviderSet is data providers.
var ServiceProviderSet = wire.NewSet(NewUserService)

type User struct {
	Id          int64
	Email       string
	PhoneNumber string
	PassWord    string
}
