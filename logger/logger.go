package logger

type LogLevel int

const (
	// TraceLevel is the constant to use when setting the Trace level for loggers
	// provided by this library.
	TraceLevel LogLevel = iota

	// DebugLevel is the constant to use when setting the Debug level for loggers
	// provided by this library.
	DebugLevel

	// InfoLevel is the constant to use when setting the Info level for loggers
	// provided by this library.
	InfoLevel

	// WarnLevel is the constant to use when setting the Warn level for loggers
	// provided by this library.
	WarnLevel

	// ErrorLevel is the constant to use when setting the Error level for loggers
	// provided by this library.
	ErrorLevel

	// PanicLevel is the constant to use when setting the Panic level for loggers
	// provided by this library.
	PanicLevel

	// FatalLevel is the constant to use when setting the Fatal level for loggers
	// provided by this library.
	FatalLevel
)

// Fields represents a map of key-value pairs where the value can be any Go
// type. The value must be able to be converted to a string.
type Fields map[string]any

// Logger is an interface for Logging
type Logger interface {
	// Trace logs a message at the Trace level
	Trace(msg ...any)

	// Tracef formats a message according to a format specifier and logs the
	// message at the Trace level
	Tracef(template string, args ...any)

	// Tracew logs a message at the Trace level along with some additional
	// context (key-value pairs)
	Tracew(msg string, fields Fields)

	// Debug logs a message at the Debug level
	Debug(msg ...any)

	// Debugf formats a message according to a format specifier and logs the
	// message at the Debug level
	Debugf(template string, args ...any)

	// Debugw logs a message at the Debug level along with some additional
	// context (key-value pairs)
	Debugw(msg string, fields Fields)

	// Info logs a message at the Info level
	Info(msg ...any)

	// Infof formats a message according to a format specifier and logs the
	// message at the Info level
	Infof(template string, args ...any)

	// Infow logs a message at the Info level along with some additional
	// context (key-value pairs)
	Infow(msg string, fields Fields)

	// Warn logs a message at the Warn level
	Warn(msg ...any)

	// Warnf formats a message according to a format specifier and logs the
	// message at the Warning level
	Warnf(template string, args ...any)

	// Warnw logs a message at the Warning level along with some additional
	// context (key-value pairs)
	Warnw(msg string, fields Fields)

	// Error logs a message at the Error level
	Error(msg ...any)

	// Errorf formats a message according to a format specifier and logs the
	// message at the Error level
	Errorf(template string, args ...any)

	// Errorw logs a message at the Error level along with some additional
	// context (key-value pairs)
	Errorw(msg string, fields Fields)

	// Panic logs a message at the Panic level and panics
	Panic(msg ...any)

	// Panicf formats a message according to a format specifier and logs the
	// message at the Panic level and then panics
	Panicf(template string, args ...any)

	// Panicw logs a message at the Panic level along with some additional
	// context (key-value pairs) and then panics
	Panicw(msg string, fields Fields)

	// Fatal logs a message at the Fatal level and exists the application
	Fatal(msg ...any)

	// Fatalf formats a message according to a format specifier and logs the
	// message at the Fatal level and exits the application
	Fatalf(template string, args ...any)

	// Fatalw logs a message at the Fatal level along with some additional
	// context (key-value pairs) and exits the application
	Fatalw(msg string, fields Fields)
}
