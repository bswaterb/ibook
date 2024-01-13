package data

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ibook/internal/service"
	"time"
)

var (
	//go:embed script/lua/article/interactive_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type articleInteractiveRepo struct {
	data *Data
}

type articleInteractiveCache struct {
	data *Data
}

func NewArticleInteractiveRepo(data *Data) service.ArticleInteractiveRepo {
	return &articleInteractiveRepo{data: data}
}

func newArticleInteractiveCache(data *Data) service.ArticleInteractiveCache {
	return &articleInteractiveCache{data: data}
}

func (repo *articleInteractiveRepo) IncrReadCount(ctx *gin.Context, articleId int64) error {
	now := time.Now().UTC().UnixMilli()
	return repo.data.mdb.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt":     gorm.Expr("read_cnt + 1"),
			"updated_time": now,
		}),
	}).Create(&Interactive{
		ArticleId:   articleId,
		ReadCnt:     1,
		CreatedTime: now,
		UpdatedTime: now,
	}).Error
}

func (repo *articleInteractiveRepo) IncrLikeCnt(ctx *gin.Context, articleId int64) error {
	now := time.Now().UTC().UnixMilli()
	return repo.data.mdb.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"like_cnt":     gorm.Expr("like_cnt + 1"),
			"updated_time": now,
		}),
	}).Create(&Interactive{
		ArticleId:   articleId,
		LikeCnt:     1,
		CreatedTime: now,
		UpdatedTime: now,
	}).Error
}

func (repo *articleInteractiveRepo) DecrLikeCnt(ctx *gin.Context, articleId int64) error {
	now := time.Now().UTC().UnixMilli()
	return repo.data.mdb.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"like_cnt":     gorm.Expr("like_cnt - 1"),
			"updated_time": now,
		}),
	}).Create(&Interactive{
		ArticleId:   articleId,
		LikeCnt:     0,
		CreatedTime: now,
		UpdatedTime: now,
	}).Error
}

func (cache *articleInteractiveCache) IncrReadCountInCache(ctx *gin.Context, articleId int64) error {
	return cache.data.rdb.Eval(ctx, luaIncrCnt, []string{genArticleInteractiveCacheKey(articleId)}, fieldReadCnt, 1).Err()
}

func (cache *articleInteractiveCache) IncrLikeCountInCache(ctx *gin.Context, articleId int64) error {
	return cache.data.rdb.Eval(ctx, luaIncrCnt, []string{genArticleInteractiveCacheKey(articleId)}, fieldLikeCnt, 1).Err()
}

func (cache *articleInteractiveCache) DecrLikeCountInCache(ctx *gin.Context, articleId int64) error {
	return cache.data.rdb.Eval(ctx, luaIncrCnt, []string{genArticleInteractiveCacheKey(articleId)}, fieldLikeCnt, -1).Err()
}

func genArticleInteractiveCacheKey(articleId int64) string {
	return fmt.Sprintf("article:interactive:%d", articleId)
}
