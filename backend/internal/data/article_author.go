package data

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"time"
)

type ArticleAuthor struct {
	Article
}

type articleAuthorRepo struct {
	db *Data
}

func (repo *articleAuthorRepo) CreateArticle(ctx *gin.Context, article *service.ArticleAuthor) (bool, error) {
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

func (repo *articleAuthorRepo) UpdateArticle(ctx *gin.Context, article *service.ArticleAuthor) (bool, error) {
	now := time.Now().UTC().UnixMilli()
	newArticle := &ArticleAuthor{
		Article{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			UpdatedTime: now,
		},
	}

	// 不更新 author_id 以及 created_time 字段
	res := repo.db.mdb.Model(newArticle).Omit("created_time", "author_id").Where("author_id=?", article.Author.Id).Save(newArticle)
	if err := res.Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewArticleAuthorRepo(db *Data) service.ArticleAuthorRepo {
	return &articleAuthorRepo{db: db}
}
