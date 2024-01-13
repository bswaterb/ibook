package data

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ibook/internal/service"
	"ibook/pkg/utils/logger"
	"time"
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

type Interactive struct {
	Id          int64 `gorm:"primaryKey,autoIncrement"`
	ArticleId   int64 `gorm:"unique"`
	ReadCnt     int64 // 阅读数
	LikeCnt     int64 // 点赞数
	CollectCnt  int64 // 收藏数
	CreatedTime int64
	UpdatedTime int64
}

type LikeRecord struct {
	Id          int64 `gorm:"primaryKey,autoIncrement"`
	UserId      int64 `gorm:"index=uid_aid, unique"`
	ArticleId   int64 `gorm:"index=uid_aid, unique"`
	CreatedTime int64
	UpdatedTime int64
	Status      uint8
}

type articleSyncRepo struct {
	data   *Data
	logger logger.Logger
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
		readerRepo := NewArticleReaderRepo(repo.data, repo.logger)
		return readerRepo.UpsertArticle(ctx, articleR)
	})
	return err
}

func (repo *articleSyncRepo) SyncUpdateStatus(ctx *gin.Context, articleId int64, authorId int64, status uint8) error {
	err := repo.data.mdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		mdbData := &Data{mdb: tx}
		authorRepo := NewArticleAuthorRepo(mdbData)
		readerRepo := NewArticleReaderRepo(mdbData, repo.logger)
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

func (repo *articleSyncRepo) UpsertLikeInfo(ctx *gin.Context, userId int64, articleId int64) error {
	now := time.Now().UTC().UnixMilli()
	err := repo.data.mdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 向记录表中插入记录
		er := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"updated_time": now,
				"status":       1,
			}),
		}).Create(&LikeRecord{
			UserId:      userId,
			ArticleId:   articleId,
			CreatedTime: now,
			UpdatedTime: now,
			Status:      1,
		}).Error
		if er != nil {
			return er
		}
		return tx.Clauses(clause.OnConflict{
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
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *articleSyncRepo) CancelLikeInfo(ctx *gin.Context, userId int64, articleId int64) error {
	now := time.Now().UTC().UnixMilli()
	err := repo.data.mdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将记录隐藏
		res := tx.Model(&LikeRecord{}).
			Where("user_id=? and articleId=?", userId, articleId).
			Updates(map[string]any{
				"status":       0,
				"updated_time": now,
			})
		if res.Error != nil {
			return nil
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("未找到相应点赞记录")
		}
		return tx.Clauses(clause.OnConflict{
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
	})
	return err
}

func NewArticleSyncRepo(data *Data, mylogger logger.Logger) service.ArticleSyncRepo {
	return &articleSyncRepo{data: data, logger: mylogger}
}
