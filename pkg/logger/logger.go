package logger

import "go.uber.org/zap"

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
