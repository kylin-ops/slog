package slog

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

func NewLfsHook(writer io.Writer, formatter logrus.Formatter) *lfshook.LfsHook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatter)
	return lfsHook
}

func setLogLevel(level string) logrus.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warm":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "trace":
		return logrus.TraceLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		logrus.Warn("日志级别设置错误，使用默认日志级别:\"info\"")
		return logrus.InfoLevel
	}
}
