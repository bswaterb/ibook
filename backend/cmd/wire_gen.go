// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/conf"
	"ibook/internal/data"
	"ibook/internal/data/message/sms/mem"
	"ibook/internal/service"
	"ibook/internal/web"
)

// Injectors from wire.go:

// wireApp init gin application.
func wireApp(secret *conf.Secret, mySQL *conf.MySQL, redis *conf.Redis, server *conf.Server) (*gin.Engine, func(), error) {
	db := data.NewMDB(mySQL)
	cmdable := data.NewRDB(redis)
	dataData, cleanup := data.NewData(db, cmdable)
	userRepo := data.NewUserRepo(dataData)
	userCache := data.NewUserCache(cmdable)
	smsRepo := mem.NewMemSMSRepo()
	verifyCodeRepo := data.NewVerifyCodeRepo(cmdable)
	userService := service.NewUserService(userRepo, userCache, smsRepo, verifyCodeRepo)
	userHandler := web.NewUserHandler(userService)
	v := newMiddleware(secret, cmdable)
	engine := newApp(userHandler, v)
	return engine, func() {
		cleanup()
	}, nil
}
