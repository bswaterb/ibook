package data

import (
	"context"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ibook/internal/conf"
)

// DataProviderSet is data providers.
var DataProviderSet = wire.NewSet(NewData, NewUserRepo)

type Data struct {
	rdb *redis.Client
	mdb *gorm.DB
}

func NewData(dataConf *conf.Data) *Data {
	return &Data{
		rdb: initRDB(dataConf.RedisConf),
		mdb: initMDB(dataConf.MysqlConf),
	}
}

func initMDB(mConf *conf.MySQL) *gorm.DB {
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

func initRDB(rConf *conf.Redis) *redis.Client {
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
	return db.AutoMigrate(&User{})
}
