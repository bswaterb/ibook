package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
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
	server, cleanup, err := wireApp(config.SecretConf, config.DataConf, config.ServerConf)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := server.Run(config.ServerConf.Port); err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.GetConf().DataConf.RedisConf.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, conf.GetConf().ServerConf.Domain)
		},
		MaxAge: 12 * time.Hour,
	}))

	store := memstore.NewStore([]byte("bswaterb1234567"),
		[]byte("bswaterb7654321"))

	server.Use(sessions.Sessions("mysession", store))
	jwtauth.SetEncryptEnv(conf.GetConf().SecretConf.JwtConf.Key, conf.GetConf().SecretConf.JwtConf.LifeDurationTime)
	server.Use(jwtauth.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	return server
}

func newApp(userHandler *web.UserHandler) *gin.Engine {
	sever := gin.Default()
	ug := sever.Group("/user")
	userHandler.RegisterRoutesV1(ug)
	return sever
}
