package client

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type LoggerBuilder struct {
	opts *options.LoggerOptions
}

func NewLoggerBuilder() *LoggerBuilder {
	return &LoggerBuilder{opts: options.Logger()}
}

// ComponentLevel sets the LogLevel value for a LogComponent.
func (l *LoggerBuilder) ComponentLevel(component options.LogComponent, level options.LogLevel) *LoggerBuilder {
	l.opts.SetComponentLevel(component, level)
	return l
}

// MaxDocumentLength sets the maximum length of a document to be logged.
func (l *LoggerBuilder) MaxDocumentLength(maxDocumentLength uint) *LoggerBuilder {
	l.opts.SetMaxDocumentLength(maxDocumentLength)
	return l
}

// Sink sets the LogSink to use for logging.
func (l *LoggerBuilder) Sink(sink options.LogSink) *LoggerBuilder {
	l.opts.SetSink(sink)
	return l
}

func (l *LoggerBuilder) build() *options.LoggerOptions {
	return l.opts
}

type slogSink struct{}

func SlogSink() options.LogSink {
	return &slogSink{}
}

func (s *slogSink) Info(level int, message string, keysAndValues ...interface{}) {
	switch options.LogLevel(level) {
	case options.LogLevelInfo:
		slog.Info(message, keysAndValues...)
	case options.LogLevelDebug:
		slog.Info(message, keysAndValues...)
	}
}

func (s *slogSink) Error(err error, message string, keysAndValues ...interface{}) {
	slog.Error(message, append(keysAndValues, "error", err)...)
}
