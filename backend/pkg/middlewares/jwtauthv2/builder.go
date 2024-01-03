package jwtauthv2

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"ibook/internal/conf"
	"net/http"
	"strings"
	"time"
)

// 实现长短 token 判断登录
var (
	shortSecretKey                 = "mysecretkey1"
	shortLifeDurationSeconds int64 = 60 * 60
	longSecretKey                  = "mysecretkey2"
	longLifeDurationSeconds  int64 = 60 * 60 * 24 * 30
)

type shortTokenClaims struct {
	jwt.RegisteredClaims
	userId  int64
	tokenId string
}

type longTokenClaims struct {
	jwt.RegisteredClaims
	userId  int64
	tokenId string
}

// GenerateTokens 生成长短 token [0]-> short [1] -> long [2] -> tokenUUID
func GenerateTokens(ctx *gin.Context, userId int64) []string {
	now := time.Now().UTC()
	tokenId := uuid.NewString()
	lifeTime := now.Add(time.Duration(int64(time.Second) * shortLifeDurationSeconds))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, shortTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifeTime),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		userId:  userId,
		tokenId: tokenId,
	})
	shortToken, err := token.SignedString([]byte(shortSecretKey))
	if err != nil {
		panic(err)
	}
	lifeTime = now.Add(time.Duration(int64(time.Second) * longLifeDurationSeconds))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, longTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifeTime),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		userId:  userId,
		tokenId: tokenId,
	})
	longToken, err := token.SignedString([]byte(shortSecretKey))
	if err != nil {
		panic(err)
	}
	return []string{shortToken, longToken, tokenId}
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

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
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
		lsTokens := strings.Split(tokenStr, "&&")
		if len(lsTokens) != 2 {
			unAuthorized(ctx)
			return
		}
		shortTokenStr := lsTokens[0]
		// 校验短 Token 是否过期
		shortClaims := &shortTokenClaims{}
		shortToken, err := jwt.ParseWithClaims(shortTokenStr, shortClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.GetConf().SecretConf.JwtConf.Key), nil
		})
		if err != nil {
			unAuthorized(ctx)
			return
		}
		if shortToken == nil {
			unAuthorized(ctx)
			return
		}
		// 短 Token 已过期，要求用户使用长 Token 进行刷新
		if !shortToken.Valid {
			shortTokenExpired(ctx)
			return
		}
		ctx.Set("userId", shortClaims.userId)
	}
}

func unAuthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "未认证无法访问此接口",
	})
	ctx.Abort()
}

func shortTokenExpired(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "用户状态暂过期，请刷新后重试",
	})
	ctx.Abort()
}
