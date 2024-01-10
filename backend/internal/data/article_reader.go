package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ibook/internal/service"
	"ibook/pkg/utils/bskit/slice"
	"ibook/pkg/utils/logger"
	"time"
)

var (
	NotInCacheErr   = fmt.Errorf("缓存中无此数据")
	CacheUnknownErr = fmt.Errorf("缓存服务异常")
)

type ArticleReader struct {
	Article
}

type articleReaderRepo struct {
	db     *Data
	logger logger.Logger
}

func (repo *articleReaderRepo) UpsertArticle(ctx *gin.Context, article *service.ArticleReader) error {
	now := time.Now().UnixMilli()
	newArticle := &ArticleReader{
		Article{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			AuthorId:    article.Author.Id,
			CreatedTime: now,
			UpdatedTime: now,
			Status:      uint8(article.Status),
		},
	}
	err := repo.db.mdb.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":        newArticle.Title,
			"content":      newArticle.Content,
			"updated_time": now,
			"status":       newArticle.Status,
		}),
	}).Create(newArticle).Error
	return err
}

func (repo *articleReaderRepo) UpdateStatusById(ctx *gin.Context, articleId int64, authorId int64, status uint8) error {
	res := repo.db.mdb.WithContext(ctx).Model(&ArticleReader{}).
		Where("id=? and author_id=?", articleId, authorId).
		Updates(map[string]any{
			"status": status,
		})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("文章已隐藏或作者不对应")
	}
	return nil
}

func (repo *articleReaderRepo) ListAll(ctx *gin.Context, offset int64, limit int64) ([]*service.Article, error) {
	l := logger.TagCtxLogger(ctx, repo.logger, "articleReaderRepo - ListAll")
	cacheRes, err := repo.listArticlesInCache(ctx, 0)
	// 缓存查询成功
	if err == nil {
		return cacheRes, nil
	}
	// 缓存查询失败，且缓存服务异常
	if !errors.Is(err, NotInCacheErr) {
		// 考虑降级
		l.Error("缓存服务不可用", logger.Field{
			Key:   "详情",
			Value: err,
		})
		return nil, fmt.Errorf("缓存服务异常，触发降级")
	}
	// 缓存中无数据，需要进入 DB 查询
	var dbRes []Article
	status := repo.db.mdb.WithContext(ctx).Model(&Article{}).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("updated_time desc").Find(&dbRes)
	if status.Error != nil {
		l.Error("MySQL 查询时出错", logger.Field{
			Key:   "详情",
			Value: err,
		})
		return nil, err
	}
	if status.RowsAffected == 0 {
		return []*service.Article{}, nil
	}
	res := slice.Map[Article, *service.Article](dbRes, func(idx int, src Article) *service.Article {
		return &service.Article{
			Id:       src.Id,
			Title:    src.Title,
			Abstract: "",
			Content:  src.Content,
			Status:   service.ArticleStatus(src.Status),
			Author: service.Author{
				Id:   src.AuthorId,
				Name: "功能未完善",
			},
			UpdatedTime: src.UpdatedTime,
			CreatedTime: src.CreatedTime,
		}
	})
	go func() {
		// 异步将结果写入缓存
		bs, err := json.Marshal(res)
		if err != nil {
			// 序列化失败
			l.Error("json序列化失败")
			return
		}
		_, err = repo.db.rdb.Set(ctx, "article:first_page:0", bs, time.Minute*10).Result()
		if err != nil {
			// 回写缓存失败
			l.Error("缓存回写失败", logger.Field{
				Key:   "详情",
				Value: err,
			})
		}
		// 预写缓存
	}()
	return res, nil
}

func (repo *articleReaderRepo) ListById(ctx *gin.Context, userId int64, offset int64, limit int64) ([]*service.Article, error) {
	panic("implement me")
}

func NewArticleReaderRepo(db *Data, myLogger logger.Logger) service.ArticleReaderRepo {
	return &articleReaderRepo{db: db, logger: myLogger}
}

func firstPageKey(author int64) string {
	return fmt.Sprintf("article:first_page:%d", author)
}

func (repo *articleReaderRepo) listArticlesInCache(ctx *gin.Context, userId int64) ([]*service.Article, error) {
	l := logger.TagCtxLogger(ctx, repo.logger, "articleReaderRepo - listArticlesInCache")
	// 先查缓存
	cacheData, err := repo.db.rdb.Get(ctx, firstPageKey(userId)).Bytes()
	if err != nil {
		// 考虑降级
		l.Warn("文章Redis缓存查询失败", logger.Field{
			Key:   "详情",
			Value: err,
		})
		return nil, err
	}
	var articles []*service.Article
	err = json.Unmarshal(cacheData, &articles)
	// 查询成功
	if len(articles) > 0 {
		return articles, nil
	}
	return nil, NotInCacheErr
}
