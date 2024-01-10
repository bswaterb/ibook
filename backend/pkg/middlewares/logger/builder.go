package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ibook/pkg/utils/logger"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	logger        logger.Logger
}

func NewBuilder(logger logger.Logger) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logger: logger,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			// URL 本身也可能很长
			Url: url,
		}
		reqId := uuid.NewString()
		specLogger := b.logger.With(logger.Field{
			Key:   "request-uuid",
			Value: reqId,
		})
		ctx.Set("ctx-logger", specLogger)
		specLogger.Info("[请求入参]:", []logger.Field{{"Method", al.Method}, {"URL", al.Url}}...)
		if b.allowReqBody && ctx.Request.Body != nil {
			// Body 读完就没有了
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			specLogger.Info("[请求出参]:", transferLogStructToString(al)...)
		}()
		// 执行业务逻辑
		ctx.Next()

	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求 URL
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}

func transferLogStructToString(log *AccessLog) []logger.Field {
	return []logger.Field{
		{
			"Duration",
			log.Duration,
		},
		{
			"Status",
			log.Status,
		},
		{
			"ResponseBody",
			log.RespBody,
		},
	}
}
