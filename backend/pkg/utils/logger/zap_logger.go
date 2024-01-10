package logger

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &zapLogger{logger: l}
}

func (z *zapLogger) Debug(msg string, args ...Field) {
	z.logger.Debug(msg, toZapFields(args)...)
}

func (z *zapLogger) Info(msg string, args ...Field) {
	z.logger.Info(msg, toZapFields(args)...)
}

func (z *zapLogger) Warn(msg string, args ...Field) {
	z.logger.Warn(msg, toZapFields(args)...)
}

func (z *zapLogger) Error(msg string, args ...Field) {
	z.logger.Error(msg, toZapFields(args)...)
}

func (z *zapLogger) With(args ...Field) Logger {
	l := z.logger.With(toZapFields(args)...)
	return &zapLogger{logger: l}
}

func toZapFields(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}
