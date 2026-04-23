package logger

import "go.uber.org/zap"

var log *zap.Logger

func init() {
	l, _ := zap.NewProduction()
	log = l
}

func GetLogger() *zap.Logger {
	return log
}
