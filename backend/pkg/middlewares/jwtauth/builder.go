package jwtauth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"ibook/internal/conf"
	"net/http"
	"strings"
	"time"
)

var (
	secretKey           string
	lifeDurationSeconds int64
)

type myClaims struct {
	jwt.RegisteredClaims
	userId int64
}

func GenerateToken(userId int64) string {
	now := time.Now().UTC()
	lifeTime := now.Add(time.Duration(int64(time.Second) * lifeDurationSeconds))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifeTime),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		userId: userId,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err)
	}
	return tokenString
}

type LoginJWTMiddlewareBuilder struct {
	paths               []string
	secretKey           string
	lifeDurationSeconds int64
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// SetEncryptEnv 用户需要传入用以加密的 secretKey 与 token 的有效时长（单位为秒）
func SetEncryptEnv(userSecretKey string, userLifeDurationSeconds int64) {
	secretKey = userSecretKey
	lifeDurationSeconds = userLifeDurationSeconds
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	if secretKey == "" && lifeDurationSeconds == 0 {
		panic("请在调用 Build 方法前先调用 jwtauth.SetEncryptEnv 方法设置 jwt 的必要信息")
	}
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			unAuthorized(ctx)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			unAuthorized(ctx)
			return
		}
		tokenStr := segs[1]
		claims := &myClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.GetConf().SecretConf.JwtConf.Key), nil
		})
		if err != nil {
			unAuthorized(ctx)
			return
		}
		if token == nil || !token.Valid {
			unAuthorized(ctx)
			return
		}

		now := time.Now().UTC()
		// 每十秒钟刷新一次
		if claims.IssuedAt.Sub(now) < -time.Second*10 {
			newToken := GenerateToken(claims.userId)
			ctx.Header("x-jwt-token", newToken)
		}
		ctx.Set("userId", claims.userId)
	}
}

func unAuthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "未认证无法访问此接口",
	})
	ctx.Abort()
}
