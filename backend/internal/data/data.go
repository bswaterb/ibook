package data

import (
	"context"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ibook/internal/conf"
	"ibook/internal/data/message/sms/ratelimit"
	"log"
)

// DataProviderSet is data providers.
var DataProviderSet = wire.NewSet(NewData,
	NewUserRepo, NewVerifyCodeRepo, NewArticleReaderRepo, NewArticleAuthorRepo, NewArticleSyncRepo, NewUserCache, NewMDB, NewRDB, ratelimit.NewRateLimitSmsRepo)

type Data struct {
	rdb redis.Cmdable
	mdb *gorm.DB
}

func NewData(mdb *gorm.DB, rdb redis.Cmdable) (*Data, func()) {
	cleanup := func() {
		log.Println("closing the data repo...")
	}
	return &Data{
		rdb: rdb,
		mdb: mdb,
	}, cleanup
}

func NewMDB(mConf *conf.MySQL) *gorm.DB {
	db, err := gorm.Open(mysql.Open(mConf.DSN))
	if err != nil {
		panic("初始化 MySQL 连接失败: " + err.Error())
	}
	err = initTable(db)
	if err != nil {
		panic("MySQL 自动迁移失败: " + err.Error())
	}
	return db
}

func NewRDB(rConf *conf.Redis) redis.Cmdable {
	url, err := redis.ParseURL(rConf.Addr)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(url)
	status := rdb.Ping(context.Background())
	if _, err = status.Result(); err != nil {
		panic("Redis初始化失败，检查Rdb服务状态")
	}
	return rdb
}

func initTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, ArticleAuthor{}, ArticleReader{}, Interactive{}, LikeRecord{})
}
