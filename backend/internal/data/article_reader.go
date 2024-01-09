package data

import (
	"fmt"
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

func NewArticleReaderRepo(db *Data) service.ArticleReaderRepo {
	return &articleReaderRepo{db: db}
}
