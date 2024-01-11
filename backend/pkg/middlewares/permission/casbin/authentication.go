package casbin

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v2"
	cb "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"ibook/internal/conf"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const rule = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

var enforcer *cb.Enforcer
var cache *bigcache.BigCache

func InitWithMySQL(dbconf *conf.MySQL) {
	db, err := gorm.Open("mysql", dbconf.DSN)
	if err != nil {
		panic(fmt.Errorf("casbin 鉴权中间件初始化失败：%w", err))
	}
	// mysql 适配器
	adapter := gormadapter.NewAdapterByDB(db)
	// 通过mysql适配器新建一个enforcer
	enforc, _ := cb.NewEnforcer(rule, adapter)
	// 日志记录
	enforc.EnableLog(true)
	enforcer = enforc
}

func InitCache() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Minute))
}

func Build() func(ctx *gin.Context) {
	if enforcer == nil {
		panic("在 Build 前请先调用 Init 方法初始化 Enforcer")
	}
	if cache == nil {
		panic("在 Build 前请先调用 Init 方法初始化 Policy-Cache")
	}
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists {
			unAuthorized(ctx)
		}
		reqPath := ctx.Request.URL.Path
		reqMethod := ctx.Request.Method
		cacheKey := strings.Join([]string{strconv.Itoa(userId.(int)), reqPath, reqMethod}, "-")
		entry, err := cache.Get(cacheKey)
		if err != nil && entry != nil {
			if string(entry) == "true" {
				ctx.Next()
			} else {
				noPrivilege(ctx)
			}
		} else if errors.Is(err, bigcache.ErrEntryNotFound) {
			// 进入数据库重新加载策略
			err := enforcer.LoadPolicy()
			if err != nil {
				serverError(ctx)
			}
			allow, err := enforcer.Enforce(userId, reqPath, reqMethod)
			if err != nil {
				serverError(ctx)
			}
			if allow {
				err := cache.Set(cacheKey, []byte("true"))
				if err != nil {
					serverError(ctx)
				}
				ctx.Next()
			} else {
				err := cache.Set(cacheKey, []byte("false"))
				if err != nil {
					serverError(ctx)
				}
				noPrivilege(ctx)
			}
		} else {
			serverError(ctx)
		}
	}
}

func unAuthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "未认证无法访问此接口",
	})
	ctx.Abort()
}

func noPrivilege(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusNotAcceptable,
		"msg":  "不具备访问该接口的权限",
	})
	ctx.Abort()
}

func serverError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "服务内部故障",
	})
	ctx.Abort()
}
