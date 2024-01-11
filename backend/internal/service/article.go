package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	mylogger "ibook/pkg/utils/logger"
)

var (
	ArticleUnKnownErr = errors.New("文章服务未知错误")
)

type ArticleAuthorRepo interface {
	CreateArticle(ctx *gin.Context, article *ArticleAuthor) (bool, error)
	UpdateArticle(ctx *gin.Context, article *ArticleAuthor) (bool, error)
	UpdateStatusById(ctx *gin.Context, articleId int64, authorId int64, status uint8) error
	GetArticleById(ctx *gin.Context, id int64, userId int64) (*Article, error)
}

type ArticleReaderRepo interface {
	UpsertArticle(ctx *gin.Context, article *ArticleReader) error
	UpdateStatusById(ctx *gin.Context, articleId int64, authorId int64, status uint8) error
	ListAll(ctx *gin.Context, offset int64, limit int64) ([]*Article, error)
	ListById(ctx *gin.Context, userId int64, offset int64, limit int64) ([]*Article, error)
}

type ArticleSyncRepo interface {
	Sync(ctx *gin.Context, articleA *ArticleAuthor, articleR *ArticleReader) error
	SyncUpdateStatus(ctx *gin.Context, articleId int64, authorId int64, status uint8) error
}

type ArticleService interface {
	EditArticle(ctx *gin.Context, article *Article) error
	PublishArticle(ctx *gin.Context, article *Article) error
	WithDrawArticle(ctx *gin.Context, articleId, userId int64) error
	GetArticleDetail(ctx *gin.Context, articleId int64, userId int64) (*Article, error)
	ListPubArticles(ctx *gin.Context, userId int64, offset int64, limit int64) ([]*Article, error)
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
	l := mylogger.TagCtxLogger(ctx, service.logger, "EditArticle")

	articleA := &ArticleAuthor{Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author,
	}}
	if article.Id <= 0 {
		// 创建新的
		articleA.Status = ArticleStatusUnpublished
		ok, err := service.ar.CreateArticle(ctx, articleA)
		if err != nil || !ok {
			l.Warn("新建文章时出现错误", mylogger.Field{
				Key:   "错误详情",
				Value: err,
			})
			return err
		}
		article.Id = articleA.Id
		return nil
	} else {
		// 更新已有的
		articleA.Status = ArticleStatusUnpublished
		ok, err := service.ar.UpdateArticle(ctx, articleA)
		if err != nil || !ok {
			l.Warn("更新文章时出现错误", mylogger.Field{
				Key:   "错误详情",
				Value: err,
			})
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
		Status: ArticleStatusPublished,
	}}
	articleR := &ArticleReader{Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author,
		Status:  ArticleStatusPublished,
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

func (service *articleService) WithDrawArticle(ctx *gin.Context, articleId, authorId int64) error {
	return service.sr.SyncUpdateStatus(ctx, articleId, authorId, ArticleStatusPrivate)
}

func (service *articleService) ListPubArticles(ctx *gin.Context, authorId int64, offset int64, limit int64) ([]*Article, error) {
	if authorId <= 0 {
		// 获取全部文章列表
		articles, err := service.rr.ListAll(ctx, offset, limit)
		if err != nil {
			// 出错
		}
		return articles, nil
	} else {
		// 获取指定作者 id 的文章列表
		articles, err := service.rr.ListById(ctx, authorId, offset, limit)
		if err != nil {
			// 出错
		}
		return articles, nil
	}
}

func (service *articleService) GetArticleDetail(ctx *gin.Context, articleId int64, userId int64) (*Article, error) {
	service.logger.Debug("")
	article, err := service.ar.GetArticleById(ctx, articleId, userId)
	if err != nil {
		return nil, fmt.Errorf("查询出错，文章不存在或用户不匹配")
	}
	return article, nil
}
