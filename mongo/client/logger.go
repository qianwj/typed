package client

import (
	"github.com/qianwj/typed/mongo/options"
	rawoptions "go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type LoggerBuilder struct {
	opts *rawoptions.LoggerOptions
}

func NewLoggerBuilder() *LoggerBuilder {
	return &LoggerBuilder{opts: rawoptions.Logger()}
}

// AddComponentLevel add the LogLevel value for a LogComponent.
func (l *LoggerBuilder) AddComponentLevel(component options.LogComponent, level options.LogLevel) *LoggerBuilder {
	l.opts.SetComponentLevel(rawoptions.LogComponent(component), rawoptions.LogLevel(level))
	return l
}

// MaxDocumentLength sets the maximum length of a document to be logged.
func (l *LoggerBuilder) MaxDocumentLength(maxDocumentLength uint) *LoggerBuilder {
	l.opts.SetMaxDocumentLength(maxDocumentLength)
	return l
}

// Sink sets the LogSink to use for logging.
func (l *LoggerBuilder) Sink(sink rawoptions.LogSink) *LoggerBuilder {
	l.opts.SetSink(sink)
	return l
}

func (l *LoggerBuilder) build() *rawoptions.LoggerOptions {
	return l.opts
}

type slogSink struct {
	logger *slog.Logger
}

func SlogSink(logger *slog.Logger) rawoptions.LogSink {
	return &slogSink{
		logger: logger,
	}
}

func (s *slogSink) Info(level int, message string, keysAndValues ...interface{}) {
	switch rawoptions.LogLevel(level) {
	case rawoptions.LogLevelInfo:
		s.logger.Info(message, keysAndValues...)
	case rawoptions.LogLevelDebug:
		s.logger.Info(message, keysAndValues...)
	}
}

func (s *slogSink) Error(err error, message string, keysAndValues ...interface{}) {
	s.logger.Error(message, append(keysAndValues, "error", err)...)
}
