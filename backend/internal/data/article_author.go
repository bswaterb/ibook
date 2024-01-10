package data

import (
	"fmt"
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
			Status:      uint8(article.Status),
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
			Status:      uint8(article.Status),
		},
	}

	// 不更新 author_id 以及 created_time 字段
	res := repo.db.mdb.WithContext(ctx).Model(newArticle).
		Where("id=? and author_id=?", article.Id, article.Author.Id).
		Updates(map[string]any{
			"title":        newArticle.Title,
			"content":      newArticle.Content,
			"updated_time": newArticle.UpdatedTime,
			"status":       newArticle.Status,
		})
	if res.RowsAffected == 0 {
		return false, fmt.Errorf("无法编辑此文章")
	}
	if err := res.Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *articleAuthorRepo) UpdateStatusById(ctx *gin.Context, articleId int64, authorId int64, status uint8) error {
	res := repo.db.mdb.WithContext(ctx).Model(&ArticleAuthor{}).
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

func (repo *articleAuthorRepo) GetArticleById(ctx *gin.Context, id int64, userId int64) (*service.Article, error) {
	article := &Article{}
	res := repo.db.mdb.WithContext(ctx).Model(&ArticleAuthor{}).
		Where("id=? and author_id=?", id, userId).First(article)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("查询无果，文章不存在或用户不匹配")
	}
	return &service.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Status:  service.ArticleStatus(article.Status),
		Author:  service.Author{Id: article.AuthorId},
	}, nil
}

func NewArticleAuthorRepo(db *Data) service.ArticleAuthorRepo {
	return &articleAuthorRepo{db: db}
}
