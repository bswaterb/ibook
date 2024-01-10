package logger

import (
	"github.com/gin-gonic/gin"
)

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) Logger
}

type Field struct {
	Key   string
	Value any
}

func GetLoggerFromCtx(ctx *gin.Context, outerLogger Logger) Logger {
	logger, exists := ctx.Get("ctx-logger")
	if !exists {
		return outerLogger
	}
	return logger.(Logger)
}

func TagLoggerFunc(outerLogger Logger, FuncName string) Logger {
	l := outerLogger.With([]Field{{"HappenedIn", FuncName}}...)
	return l
}

func TagCtxLogger(context *gin.Context, outerLogger Logger, FuncName string) Logger {
	ctxLogger := GetLoggerFromCtx(context, outerLogger)
	tagLogger := TagLoggerFunc(ctxLogger, FuncName)
	return tagLogger
}
