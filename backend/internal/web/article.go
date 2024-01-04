package web

import (
	"github.com/gin-gonic/gin"
	"ibook/internal/service"
	"ibook/pkg/utils/request"
	"ibook/pkg/utils/result"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

// RegisterRoutesV1 !#登录接口字段需对齐短信/手机号登录
func (handler *ArticleHandler) RegisterRoutesV1(server *gin.Engine) {
	ug := server.Group("/articles")
	ug.POST("/edit", handler.Edit)
}

func (handler *ArticleHandler) Edit(ctx *gin.Context) {
	req := &ArticleEditReq{}
	if err := request.ParseRequestBody(ctx, req); err != nil {
		result.RespWithError(ctx, result.PARAM_NOT_EQUAL_CODE, "请求传参或设置有误", nil)
		return
	}
	userId, _ := ctx.Get("userId")
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
		result.RespWithError(ctx, result.UNKNOWN_ERROR_CODE, "未知错误", nil)
	}
	result.RespWithSuccess(ctx, "操作成功", &ArticleEditReply{
		Id:       article.Id,
		Title:    article.Title,
		Content:  "",
		AuthorId: article.Author.Id,
	})
}
