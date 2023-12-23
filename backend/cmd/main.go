package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"ibook/internal/conf"
	"ibook/internal/web"
	"ibook/pkg/middlewares/jwtauth"
	"ibook/pkg/middlewares/ratelimit"
	"strings"
	"time"
)

func main() {
	config := conf.GetConf()
	server, cleanup, err := wireApp(config.SecretConf, config.DataConf.MysqlConf, config.DataConf.RedisConf, config.ServerConf)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := server.Run(config.ServerConf.Port); err != nil {
		panic(err)
	}
}

func newApp(userHandler *web.UserHandler, middlewares []gin.HandlerFunc) *gin.Engine {
	sever := gin.Default()
	sever.Use(middlewares...)
	// 注册 /users/*** 路由
	userHandler.RegisterRoutesV1(sever)
	return sever
}

func newMiddleware(secret *conf.Secret, redisCli *redis.Client) []gin.HandlerFunc {
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
	rlMw := ratelimit.NewBuilder(redisCli, time.Second, 100).Build()
	jwtauth.SetEncryptEnv(secret.JwtConf.Key, secret.JwtConf.LifeDurationTime)
	jwtMw := jwtauth.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		Build()
	return []gin.HandlerFunc{corsMw, rlMw, jwtMw}
}
