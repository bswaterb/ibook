package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"ibook/internal/conf"
	"ibook/internal/web"
	"ibook/pkg/middlewares/jwtauth"
	logger2 "ibook/pkg/middlewares/logger"
	"ibook/pkg/middlewares/ratelimit"
	"ibook/pkg/utils/logger"
	"strings"
	"time"
)

func main() {
	config := conf.GetConf()
	// initRemoteViper()
	fmt.Println(viper.Get("server.port"))
	server, cleanup, err := wireApp(config.SecretConf, config.DataConf.MysqlConf, config.DataConf.RedisConf, config.ServerConf)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := server.Run(config.ServerConf.Port); err != nil {
		panic(err)
	}
}

func newApp(userHandler *web.UserHandler, articleHandlers *web.ArticleHandler, middlewares []gin.HandlerFunc) *gin.Engine {
	sever := gin.Default()
	sever.Use(middlewares...)
	// 注册 /users/*** 路由
	userHandler.RegisterRoutesV1(sever)
	articleHandlers.RegisterRoutesV1(sever)
	return sever
}

func newMiddleware(secret *conf.Secret, logger logger.Logger, redisCli redis.Cmdable) []gin.HandlerFunc {
	corsMw := cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
	loggerMw := logger2.NewBuilder(logger).AllowReqBody().AllowRespBody().Build()
	rlMw := ratelimit.NewBuilder(redisCli, time.Second, 100).Build()
	jwtauth.SetEncryptEnv(secret.JwtConf.Key, secret.JwtConf.LifeDurationTime)
	jwtMw := jwtauth.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login_sms/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/users/login").
		Build()
	return []gin.HandlerFunc{corsMw, loggerMw, rlMw, jwtMw}
}

func initViper() {
	filePath := pflag.String("config", "./configs/config.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initRemoteViper() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/ibook")
	if err != nil {
		panic(err)
	}
	if err = viper.WatchRemoteConfig(); err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {

	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {

}
