package log

import (
	"context"
	"fmt"
	"io"

	logSyslog "log/syslog"

	"github.com/sirupsen/logrus/hooks/syslog"
	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go"

	"github.com/abyss414/house/app/common/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *logrus.Logger

func InitLogger() error {
	return initLogger(config.GlobalConfig().Log)
}

func initLogger(conf *config.LogConfig) error {
	switch conf.Method {
	case "file":
		return initLoggerWithFile(conf.FilePath)
	case "syslog":
		return initLoggerSyslog(conf)
	}
	return fmt.Errorf("不支持的日志输出类型%s", conf.Method)
}

func initLoggerSyslog(conf *config.LogConfig) error {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000"})
	hook, err := syslog.NewSyslogHook(conf.Syslog.Network, conf.Syslog.Addr, logSyslog.LOG_INFO, "")
	if err != nil {
		return err
	}
	logrusLogger.AddHook(hook)
	globalLogger = logrusLogger
	return nil
}

func initLoggerWithFile(fileName string) error {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000"})

	rotateFile := &lumberjack.Logger{
		Filename:   fileName,
		MaxBackups: 10,
		MaxAge:     7,
	}
	logrusLogger.SetOutput(io.MultiWriter([]io.Writer{rotateFile}...))
	globalLogger = logrusLogger
	return nil
}

func GlobalLogger() *logrus.Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

func WithContext(ctx context.Context) *logrus.Entry {
	return GlobalLogger().WithFields(transferDefaultFieldsFromContext(ctx))
}

var registeredContextKey = []string{
	"platform",
	"method",
}

func RegisterContextKey(key string) {
	registeredContextKey = append(registeredContextKey, key) //TODO Lock
}

func transferDefaultFieldsFromContext(ctx context.Context) logrus.Fields {
	result := map[string]interface{}{}
	for _, s := range registeredContextKey {
		if ctx.Value(s) != nil {
			result[s] = ctx.Value(s)
		}
	}
	// trace
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return result
	}
	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		result["trace_id"] = spanContext.TraceID().String()
		result["span_id"] = spanContext.SpanID().String()
	}
	return result
}
