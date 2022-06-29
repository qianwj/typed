package logger

import (
	"fmt"
	"log"
)

type StdLogger struct {
	Level LogLevel
}

// Trace logs a message at the Trace level
func (l StdLogger) Trace(msg ...any) {
	if l.Level <= TraceLevel {
		log.Print(append([]any{"[TRACE]   "}, msg...)...)
	}
}

// Tracef formats a message according to a format specifier and logs the
// message at the Trace level
func (l StdLogger) Tracef(template string, args ...any) {
	if l.Level <= TraceLevel {
		log.Printf("[TRACE]   "+template, args...)
	}
}

// Tracew logs a message at the Trace level along with some additional
// context (key-value pairs)
func (l StdLogger) Tracew(msg string, fields Fields) {
	if l.Level <= TraceLevel {
		log.Printf("[TRACE]   %s %s", msg, handlFields(fields))
	}
}

// Debug logs a message at the Debug level
func (l StdLogger) Debug(msg ...any) {
	if l.Level <= DebugLevel {
		log.Print(append([]any{"[DEBUG]   "}, msg...)...)
	}
}

// Debugf formats a message according to a format specifier and logs the
// message at the Debug level
func (l StdLogger) Debugf(template string, args ...any) {
	if l.Level <= DebugLevel {
		log.Printf("[DEBUG]   "+template, args...)
	}
}

// Debugw logs a message at the Debug level along with some additional
// context (key-value pairs)
func (l StdLogger) Debugw(msg string, fields Fields) {
	if l.Level <= DebugLevel {
		log.Printf("[DEBUG]   %s %s", msg, handlFields(fields))
	}
}

// Info logs a message at the Info level
func (l StdLogger) Info(msg ...any) {
	if l.Level <= InfoLevel {
		log.Print(append([]any{"[INFO]    "}, msg...)...)
	}
}

// Infof formats a message according to a format specifier and logs the
// message at the Info level
func (l StdLogger) Infof(template string, args ...any) {
	if l.Level <= InfoLevel {
		log.Printf("[INFO]    "+template, args...)
	}
}

// Infow logs a message at the Info level along with some additional
// context (key-value pairs)
func (l StdLogger) Infow(msg string, fields Fields) {
	if l.Level <= InfoLevel {
		log.Printf("[INFO]    %s %s", msg, handlFields(fields))
	}
}

// Warn logs a message at the Warn level
func (l StdLogger) Warn(msg ...any) {
	if l.Level <= WarnLevel {
		log.Print(append([]any{"[WARNING] "}, msg...)...)
	}
}

// Warnf formats a message according to a format specifier and logs the
// message at the Warning level
func (l StdLogger) Warnf(template string, args ...any) {
	if l.Level <= WarnLevel {
		log.Printf("[WARNING] "+template, args...)
	}
}

// Warnw logs a message at the Warning level along with some additional
// context (key-value pairs)
func (l StdLogger) Warnw(msg string, fields Fields) {
	if l.Level <= WarnLevel {
		log.Printf("[WARNING] %s %s", msg, handlFields(fields))
	}
}

// Error logs a message at the Error level
func (l StdLogger) Error(msg ...any) {
	if l.Level <= ErrorLevel {
		log.Print(append([]any{"[ERROR]   "}, msg...)...)
	}
}

// Errorf formats a message according to a format specifier and logs the
// message at the Error level
func (l StdLogger) Errorf(template string, args ...any) {
	if l.Level <= ErrorLevel {
		log.Printf("[ERROR]   "+template, args...)
	}
}

// Errorw logs a message at the Error level along with some additional
// context (key-value pairs)
func (l StdLogger) Errorw(msg string, fields Fields) {
	if l.Level <= ErrorLevel {
		log.Printf("[ERROR]   %s %s", msg, handlFields(fields))
	}
}

// Panic logs a message at the Panic level and panics
func (l StdLogger) Panic(msg ...any) {
	if l.Level <= PanicLevel {
		log.Panic(append([]any{"[PANIC]   "}, msg...)...)
	}
}

// Panicf formats a message according to a format specifier and logs the
// message at the Panic level and then panics
func (l StdLogger) Panicf(template string, args ...any) {
	if l.Level <= PanicLevel {
		log.Panicf("[PANIC]   "+template, args...)
	}
}

// Panicw logs a message at the Panic level along with some additional
// context (key-value pairs) and then panics
func (l StdLogger) Panicw(msg string, fields Fields) {
	if l.Level <= PanicLevel {
		log.Panicf("[PANIC]   %s %s", msg, handlFields(fields))
	}
}

// Fatal logs a message at the Fatal level and exists the application
func (l StdLogger) Fatal(msg ...any) {
	if l.Level <= FatalLevel {
		log.Fatal(append([]any{"[FATAL]   "}, msg...)...)
	}
}

// Fatalf formats a message according to a format specifier and logs the
// message at the Fatal level and exits the application
func (l StdLogger) Fatalf(template string, args ...any) {
	if l.Level <= FatalLevel {
		log.Fatalf("[FATAL]   "+template, args...)
	}
}

// Fatalw logs a message at the Fatal level along with some additional
// context (key-value pairs) and exits the application
func (l StdLogger) Fatalw(msg string, fields Fields) {
	if l.Level <= FatalLevel {
		log.Fatalf("[FATAL]   %s %s", msg, handlFields(fields))
	}
}

func handlFields(flds Fields) string {
	var ret string
	for k, v := range flds {
		ret += fmt.Sprintf("[%s=%s]", k, v)
	}
	return ret
}
