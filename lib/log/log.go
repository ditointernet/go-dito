package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Level indicates the severity of the data being logged
type Level int

const (
	// LevelCritical alerts about severe problems. Most of the time, needs some human intervention ASAP
	LevelCritical Level = iota + 1
	// LevelError alerts about events that are likely to cause problems
	LevelError
	// LevelWarning warns about events the might cause problems to the system
	LevelWarning
	// LevelInfo are routine information
	LevelInfo
	// LevelDebug are debug or trace information
	LevelDebug
)

var levelStringValueMap = map[string]Level{
	"CRITICAL": LevelCritical,
	"ERROR":    LevelError,
	"WARNING":  LevelWarning,
	"INFO":     LevelInfo,
	"DEBUG":    LevelDebug,
}

// String returns the name of the LogLevel
func (l Level) String() string {
	return []string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"INFO",
		"DEBUG",
	}[l-1]
}

// LogAttribute represents an information to be extracted from the context and included into the log
type LogAttribute string

// LogAttributeSet is a set of LogAttributes
type LogAttributeSet map[LogAttribute]bool

// LoggerInput defines the dependencies of a Logger
type LoggerInput struct {
	Level      string
	Attributes LogAttributeSet
}

// Logger is the structure responsible for log data
type Logger struct {
	level      Level
	attributes LogAttributeSet
	now        func() time.Time
}

// NewLogger constructs a new Logger instance
func NewLogger(in LoggerInput) *Logger {
	logger := &Logger{level: levelStringValueMap[in.Level], attributes: in.Attributes, now: time.Now}

	if logger.level < LevelCritical || logger.level > LevelDebug {
		logger.level = LevelInfo
	}

	return logger
}

// Debug logs debug data
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelDebug {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelDebug)
	}
}

// Info logs info data
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelInfo {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelInfo)
	}
}

// Warning logs warning data
func (l Logger) Warning(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelWarning {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelWarning)
	}
}

// Error logs error data
func (l Logger) Error(ctx context.Context, err error) {
	if l.level >= LevelError {
		l.print(ctx, err.Error(), LevelError)
	}
}

// Critical logs critical data
func (l Logger) Critical(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.print(ctx, err.Error(), LevelCritical)
	}
}

type logData struct {
	Timestamp  string                       `json:"timestamp"`
	Level      string                       `json:"level"`
	Message    string                       `json:"message"`
	Attributes map[LogAttribute]interface{} `json:"attributes,omitempty"`
}

func (l Logger) print(ctx context.Context, msg string, level Level) {
	data, _ := json.Marshal(logData{
		Timestamp:  l.now().Format(time.RFC3339),
		Level:      level.String(),
		Message:    msg,
		Attributes: l.extractLogAttributesFromContext(ctx),
	})

	fmt.Println(string(data))
}

func (l Logger) extractLogAttributesFromContext(ctx context.Context) map[LogAttribute]interface{} {
	attributes := map[LogAttribute]interface{}{}

	for attr := range l.attributes {
		if value := ctx.Value(string(attr)); value != nil {
			attributes[attr] = value
		}
	}

	return attributes
}
