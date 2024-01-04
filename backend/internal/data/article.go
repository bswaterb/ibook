package data

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"time"
)

type Article struct {
	Id          int64  `gorm:"primaryKey,authIncrement"`
	Title       string `gorm:"type=varchar(1024)"`
	Content     string `gorm:"type=blob"`
	AuthorId    int64  `gorm:"index=aid_ctime"`
	CreatedTime int64  `gorm:"index=aid_ctime"`
	UpdatedTime int64
}

type articleRepo struct {
	db *Data
}

func (repo *articleRepo) CreateArticle(ctx *gin.Context, article *service.Article) (bool, error) {
	now := time.Now().UTC().UnixMilli()
	newArticle := &Article{
		Title:       article.Title,
		Content:     article.Title,
		AuthorId:    article.Author.Id,
		CreatedTime: now,
		UpdatedTime: now,
	}
	res := repo.db.mdb.Create(newArticle)
	if err := res.Error; err != nil {
		return false, err
	}
	article.Id = newArticle.Id
	return true, nil
}

func (repo *articleRepo) UpdateArticle(ctx *gin.Context, article *service.Article) (bool, error) {
	now := time.Now().UTC().UnixMilli()
	newArticle := &Article{
		Id:          article.Id,
		Title:       article.Title,
		Content:     article.Title,
		UpdatedTime: now,
	}
	res := repo.db.mdb.Model(newArticle).Omit("created_time", "author_id").Save(newArticle)
	if err := res.Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewArticleRepo(db *Data) service.ArticleRepo {
	return &articleRepo{db: db}
}
