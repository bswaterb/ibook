// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/conf"
	"ibook/internal/data"
	ratelimit2 "ibook/internal/data/message/sms/ratelimit"
	"ibook/internal/service"
	"ibook/internal/web"
	"ibook/pkg/utils/logger"
	"ibook/pkg/utils/ratelimit"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

// wireApp init gin application.
func wireApp(secret *conf.Secret, mySQL *conf.MySQL, redis *conf.Redis, server *conf.Server) (*gin.Engine, func(), error) {
	db := data.NewMDB(mySQL)
	cmdable := data.NewRDB(redis)
	dataData, cleanup := data.NewData(db, cmdable)
	userRepo := data.NewUserRepo(dataData)
	userCache := data.NewUserCache(cmdable)
	limiter := ratelimit.NewRedisSlidingWindowLimiter(cmdable)
	smsRepo := ratelimit2.NewRateLimitSmsRepo(limiter)
	verifyCodeRepo := data.NewVerifyCodeRepo(cmdable)
	userService := service.NewUserService(userRepo, userCache, smsRepo, verifyCodeRepo)
	userHandler := web.NewUserHandler(userService)
	articleAuthorRepo := data.NewArticleAuthorRepo(dataData)
	articleReaderRepo := data.NewArticleReaderRepo(dataData)
	articleSyncRepo := data.NewArticleSyncRepo(dataData)
	loggerLogger := logger.NewZapLogger()
	articleService := service.NewArticleService(articleAuthorRepo, articleReaderRepo, articleSyncRepo, loggerLogger)
	articleHandler := web.NewArticleHandler(articleService, loggerLogger)
	v := newMiddleware(secret, loggerLogger, cmdable)
	engine := newApp(userHandler, articleHandler, v)
	return engine, func() {
		cleanup()
	}, nil
}
