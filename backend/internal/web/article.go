package web

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"ibook/pkg/utils/bskit/slice"
	"ibook/pkg/utils/logger"
	"ibook/pkg/utils/request"
	"ibook/pkg/utils/result"
	"strconv"
	"strings"
	"time"
)

const handlerName = "ArticleHandler"

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
	ug.GET("/pub/list", handler.PubList)
	ug.GET("/pub/detail/:id", handler.PubDetail)
	ug.GET("/detail/:id", handler.Detail)
	ug.POST("/pub/like", handler.Like)
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

func (handler *ArticleHandler) PubList(ctx *gin.Context) {
	l := logger.TagCtxLogger(ctx, handler.logger, "ArticleHandler-List")
	req := &ArticleListReq{}
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists || userId.(int64) < 0 {
		l.Error("未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}
	articles, err := handler.svc.ListPubArticles(ctx, req.AuthorId, req.Offset, req.Limit)
	if err != nil {
		l.Warn("获取文章列表失败", logger.Field{
			Key:   "错误详情",
			Value: err,
		})
	}
	result.RespWithSuccess(ctx, "获取成功",
		&ArticleListReply{
			Articles: slice.Map[*service.Article, *Article](articles, func(idx int, src *service.Article) *Article {
				return &Article{
					Id:          src.Id,
					Title:       src.Title,
					Abstract:    src.GenAbstract(),
					Status:      uint8(src.Status),
					AuthorId:    src.Author.Id,
					AuthorName:  src.Author.Name,
					UpdatedTime: time.Unix(src.UpdatedTime, 0).Local().Format(time.DateTime),
					CreatedTime: time.Unix(src.CreatedTime, 0).Local().Format(time.DateTime),
				}
			}),
		},
	)
}

func (handler *ArticleHandler) PubDetail(ctx *gin.Context) {
	articleId := ctx.Param("id")
	if articleId == "" {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	id, err := strconv.Atoi(articleId)
	if err != nil || id <= 0 {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	l := logger.TagCtxLogger(ctx, handler.logger, "ArticleHandler-PubDetail")
	// 从缓存中查找文章详情

	//增加计数
	go func() {
		err := handler.svc.IncrReadCount(ctx, int64(id))
		if err != nil {
			// 打日志
			l.Warn("异步增加文章阅读计数时出错：", logger.Field{
				Key:   "详情",
				Value: err,
			})
		}
	}()
}

func (handler *ArticleHandler) Detail(context *gin.Context) {
	l := logger.TagCtxLogger(context, handler.logger, "ArticleHandler-Detail")
	id := context.Param("id")
	articleId, err := strconv.ParseInt(id, 10, 64)
	if err == nil {
		result.RespWithError(context, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	userId, exists := context.Get("userId")
	if !exists || userId.(int64) < 0 {
		l.Error("未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}

	// 去制作库中获取数据
	article, err := handler.svc.GetArticleDetail(context, articleId, userId.(int64))
	if err != nil {
		l.Warn("数据库获取数据失败", logger.Field{
			Key:   "错误详情",
			Value: err,
		})
		result.RespWithError(context, result.UNKNOWN_ERROR_CODE, "文章获取失败", nil)
		return
	}
	result.RespWithSuccess(context, "获取成功", &GetArticleReply{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
	})
}

func (handler *ArticleHandler) Like(ctx *gin.Context) {
	l := logger.TagCtxLogger(ctx, handler.logger, "ArticleHandler-Like")
	var req *LikeArticleReq
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists || userId.(int64) < 0 {
		l.Error("未成功从 token 提取 userId", logger.Field{
			Key:   "token错误",
			Value: nil,
		})
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "用户未登录", nil)
	}
	// 本次标记为喜欢
	var err error
	var status string
	if req.Like == 1 {
		err = handler.svc.LikeArticle(ctx, userId.(int64), req.ArticleId)
		status = "like"
	} else {
		err = handler.svc.CancelLikeArticle(ctx, userId.(int64), req.ArticleId)
		status = "normal"
	}
	if err != nil {
		l.Warn("用户赞/取消赞文章失败", logger.Field{
			Key:   "详情",
			Value: err,
		})
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "服务内部错误", nil)
		return
	}

	result.RespWithSuccess(ctx, "操作成功", &LikeArticleReply{
		OK:            true,
		CurrentStatus: status,
	})

}
