package data

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ibook/internal/service"
	"time"
)

type ArticleReader struct {
	Article
}

type articleReaderRepo struct {
	db *Data
}

func (repo *articleReaderRepo) CreateArticle(ctx *gin.Context, article *service.ArticleReader) (bool, error) {
	now := time.Now().UTC().UnixMilli()
	newArticle := &ArticleAuthor{
		Article{
			Title:       article.Title,
			Content:     article.Content,
			AuthorId:    article.Author.Id,
			CreatedTime: now,
			UpdatedTime: now,
		},
	}
	res := repo.db.mdb.Create(newArticle)
	if err := res.Error; err != nil {
		return false, err
	}
	article.Id = newArticle.Id
	return true, nil
}

func (repo *articleReaderRepo) UpdateArticle(ctx *gin.Context, article *service.ArticleReader) (bool, error) {
	now := time.Now().UTC().UnixMilli()
	newArticle := &ArticleReader{
		Article{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			UpdatedTime: now,
		},
	}
	// 不更新 author_id 以及 created_time 字段
	res := repo.db.mdb.Model(newArticle).Omit("created_time", "author_id").Save(newArticle)
	if err := res.Error; err != nil {
		return false, err
	}
	return true, nil
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
		},
	}
	err := repo.db.mdb.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":        newArticle.Title,
			"content":      newArticle.Content,
			"updated_time": now,
		}),
	}).Create(newArticle).Error
	return err
}

func NewArticleReaderRepo(db *Data) service.ArticleReaderRepo {
	return &articleReaderRepo{db: db}
}
