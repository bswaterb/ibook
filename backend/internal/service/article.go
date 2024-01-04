package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ArticleUnKnownErr = errors.New("文章服务未知错误")
)

type ArticleRepo interface {
	CreateArticle(ctx *gin.Context, article *Article) (bool, error)
	UpdateArticle(ctx *gin.Context, article *Article) (bool, error)
}

type ArticleService interface {
	EditArticle(ctx *gin.Context, article *Article) error
}

type articleService struct {
	ar ArticleRepo
}

func NewArticleService(ar ArticleRepo) ArticleService {
	return &articleService{ar: ar}
}

func (service *articleService) EditArticle(ctx *gin.Context, article *Article) error {
	if article.Id <= 0 {
		// 创建新的
		ok, err := service.ar.CreateArticle(ctx, article)
		if err != nil || !ok {
			return err
		}
		return nil
	} else {
		// 更新已有的
		ok, err := service.ar.UpdateArticle(ctx, article)
		if err != nil || !ok {
			return err
		}
		return nil
	}
}
