package data

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ibook/internal/service"
)

type Article struct {
	Id          int64  `gorm:"primaryKey,authIncrement"`
	Title       string `gorm:"type=varchar(1024)"`
	Content     string `gorm:"type=blob"`
	AuthorId    int64  `gorm:"index=aid_ctime"`
	CreatedTime int64  `gorm:"index=aid_ctime"`
	UpdatedTime int64
	Status      uint8
}

type articleSyncRepo struct {
	data *Data
}

// Sync 同步发表文章，文章 status 都应为 published
func (repo *articleSyncRepo) Sync(ctx *gin.Context, articleA *service.ArticleAuthor, articleR *service.ArticleReader) error {
	err := repo.data.mdb.Transaction(func(tx *gorm.DB) error {
		// flag > 0 -> 更新 || flag <= 0 -> 创建
		flag := articleA.Id
		authorRepo := NewArticleAuthorRepo(repo.data)
		if flag > 0 {
			ok, err := authorRepo.UpdateArticle(ctx, articleA)
			if err != nil {
				return fmt.Errorf("同步发表过程出错：更新作者文章失败：%w", err)
			}
			if !ok {
				return fmt.Errorf("同步发表过程出错：更新作者文章失败")
			}
		} else {
			ok, err := authorRepo.CreateArticle(ctx, articleA)
			if err != nil {
				return fmt.Errorf("同步发表过程出错：创建作者文章失败：%w", err)
			}
			if !ok {
				return fmt.Errorf("同步发表过程出错：创建作者文章失败")
			}
		}
		articleR.Id = articleA.Id
		readerRepo := NewArticleReaderRepo(repo.data)
		return readerRepo.UpsertArticle(ctx, articleR)
	})
	return err
}

func (repo *articleSyncRepo) SyncUpdateStatus(ctx *gin.Context, articleId int64, authorId int64, status uint8) error {
	err := repo.data.mdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		mdbData := &Data{mdb: tx}
		authorRepo := NewArticleAuthorRepo(mdbData)
		readerRepo := NewArticleReaderRepo(mdbData)
		err := authorRepo.UpdateStatusById(ctx, articleId, authorId, status)
		if err != nil {
			return fmt.Errorf("同步更新文章状态时出错 - article_author: %w", err)
		}
		err = readerRepo.UpdateStatusById(ctx, articleId, authorId, status)
		if err != nil {
			return fmt.Errorf("同步更新文章状态时出错 - article_reader: %w", err)
		}
		return nil
	})
	return err
}

func NewArticleSyncRepo(data *Data) service.ArticleSyncRepo {
	return &articleSyncRepo{data: data}
}
