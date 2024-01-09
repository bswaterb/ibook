package web

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"ibook/pkg/utils/logger"
	"ibook/pkg/utils/request"
	"ibook/pkg/utils/result"
	"strings"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, myLogger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: myLogger,
	}
}

// RegisterRoutesV1 !#登录接口字段需对齐短信/手机号登录
func (handler *ArticleHandler) RegisterRoutesV1(server *gin.Engine) {
	ug := server.Group("/articles")
	ug.POST("/edit", handler.Edit)
	ug.POST("/publish", handler.Publish)
	ug.POST("/withdraw", handler.Withdraw)
}

func (handler *ArticleHandler) Edit(ctx *gin.Context) {
	req := &ArticleEditReq{}
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		req.Title = "未命名文章"
	}
	userId, exists := ctx.Get("userId")
	if !exists || userId.(int64) < 0 {
		handler.logger.Error("[ArticleHandler-Edit] 未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}
	article := &service.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: service.Author{
			Id: userId.(int64),
		},
	}
	err := handler.svc.EditArticle(ctx, article)
	if err != nil {
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "未知错误", err)
		return
	}
	result.RespWithSuccess(ctx, "操作成功", &ArticleEditReply{
		Id:    article.Id,
		Title: article.Title,
		// Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (handler *ArticleHandler) Publish(ctx *gin.Context) {
	req := &ArticlePublishReq{}
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		req.Title = "未命名文章"
	}
	userId, exists := ctx.Get("userId")
	if !exists || userId.(int64) < 0 {
		handler.logger.Error("[ArticleHandler-Publish] 未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}
	article := &service.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: service.Author{
			Id: userId.(int64),
		},
	}
	err := handler.svc.PublishArticle(ctx, article)
	if err != nil {
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "未知错误", nil)
		return
	}
	result.RespWithSuccess(ctx, "发布成功", &ArticlePublishReply{
		Id: article.Id,
		OK: true,
	})

}

func (handler *ArticleHandler) Withdraw(ctx *gin.Context) {
	// 将作者发表到 article_reader 表中的文章设置为隐藏状态
	req := &ArticleWithdrawReq{}
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	if req.Id <= 0 {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists || userId.(int64) < 0 {
		handler.logger.Error("[ArticleHandler-Publish] 未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}
	err := handler.svc.WithDrawArticle(ctx, req.Id, userId.(int64))
	if err != nil {
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "文章状态更新失败", nil)
		return
	}
	result.RespWithSuccess(ctx, "隐藏文章成功", &ArticleWithdrawReply{
		Id: req.Id,
		OK: true,
	})

}
