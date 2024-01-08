package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	mylogger "ibook/pkg/utils/logger"
)

var (
	ArticleUnKnownErr = errors.New("文章服务未知错误")
)

type ArticleAuthorRepo interface {
	CreateArticle(ctx *gin.Context, article *ArticleAuthor) (bool, error)
	UpdateArticle(ctx *gin.Context, article *ArticleAuthor) (bool, error)
}

type ArticleReaderRepo interface {
	CreateArticle(ctx *gin.Context, article *ArticleReader) (bool, error)
	UpdateArticle(ctx *gin.Context, article *ArticleReader) (bool, error)
	UpsertArticle(ctx *gin.Context, article *ArticleReader) error
}

type ArticleSyncRepo interface {
	Sync(ctx *gin.Context, articleA *ArticleAuthor, articleR *ArticleReader) error
}

type ArticleService interface {
	EditArticle(ctx *gin.Context, article *Article) error
	PublishArticle(ctx *gin.Context, article *Article) error
}

type articleService struct {
	ar     ArticleAuthorRepo
	rr     ArticleReaderRepo
	sr     ArticleSyncRepo
	logger mylogger.Logger
}

func NewArticleService(ar ArticleAuthorRepo, rr ArticleReaderRepo, sr ArticleSyncRepo, logger mylogger.Logger) ArticleService {
	return &articleService{ar: ar, rr: rr, sr: sr, logger: logger}
}

func (service *articleService) EditArticle(ctx *gin.Context, article *Article) error {
	articleA := &ArticleAuthor{Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author,
	}}
	if article.Id <= 0 {
		// 创建新的
		ok, err := service.ar.CreateArticle(ctx, articleA)
		if err != nil || !ok {
			return err
		}
		article.Id = articleA.Id
		return nil
	} else {
		// 更新已有的
		ok, err := service.ar.UpdateArticle(ctx, articleA)
		if err != nil || !ok {
			return err
		}
		return nil
	}
}

func (service *articleService) PublishArticle(ctx *gin.Context, article *Article) error {
	articleA := &ArticleAuthor{Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author: Author{
			Id:   article.Author.Id,
			Name: article.Author.Name,
		},
	}}
	articleR := &ArticleReader{Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author,
	}}
	successFlag := false
	retryTimes := 3
	var err error
	for !successFlag && retryTimes > 0 {
		err = service.sr.Sync(ctx, articleA, articleR)
		if err != nil {
			service.logger.Error("[ArticleService-Publish] 未成功同步双库数据：", mylogger.Field{
				Key:   "同步出错",
				Value: err,
			})
			retryTimes--
		} else {
			successFlag = true
		}
	}
	if successFlag == false {
		return err
	}
	article.Id = articleA.Id
	return nil
}
