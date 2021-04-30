package log

import (
	"context"
	"encoding/json"
	"fmt"
)

// Level indicates the severity of the data being logged
type Level string

var (
	// LevelDebug are debug or trace information
	LevelDebug Level = "DEBUG"
	// LevelInfo are routine information
	LevelInfo Level = "INFO"
	// LevelWarning warns about events the might cause problems to the system
	LevelWarning Level = "WARNING"
	// LevelError alerts about events that are likely to cause problems
	LevelError Level = "ERROR"
	// LevelCritical alerts about severe problems. Most of the time, needs some human intervention ASAP
	LevelCritical Level = "CRITICAL"
)

var levelPriority = map[Level]int{
	LevelCritical: 5,
	LevelError:    10,
	LevelWarning:  15,
	LevelInfo:     20,
	LevelDebug:    25,
}

// LogAttribute represents an information to be extracted from the context and included into the log
type LogAttribute string

// LogAttributeSet is a set og LogAttributes
type LogAttributeSet map[LogAttribute]bool

// LoggerInput defines the dependencies of a Logger
type LoggerInput struct {
	Level      Level
	Attributes LogAttributeSet
}

// Logger is the structure responsible for log data
type Logger struct {
	level      Level
	attributes LogAttributeSet
}

// NewLogger ...
func NewLogger(in LoggerInput) *Logger {
	if _, ok := levelPriority[in.Level]; !ok {
		in.Level = LevelInfo
	}

	return &Logger{level: in.Level, attributes: in.Attributes}
}

// Debug logs debug data
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if levelPriority[l.level] >= 25 {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelDebug)
	}
}

// Info logs info data
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if levelPriority[l.level] >= 20 {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelInfo)
	}
}

// Warning logs warning data
func (l Logger) Warning(ctx context.Context, msg string, args ...interface{}) {
	if levelPriority[l.level] >= 15 {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelWarning)
	}
}

// Error logs error data
func (l Logger) Error(ctx context.Context, err error) {
	if levelPriority[l.level] >= 10 {
		l.print(ctx, err.Error(), LevelError)
	}
}

// Critical logs critical data
func (l Logger) Critical(ctx context.Context, err error) {
	if levelPriority[l.level] >= 5 {
		l.print(ctx, err.Error(), LevelCritical)
	}
}

// LogData encapsulates the data being logged
type LogData struct {
	Level      Level                        `json:"level"`
	Message    string                       `json:"message"`
	Attributes map[LogAttribute]interface{} `json:"attributes,omitempty"`
}

func (l Logger) print(ctx context.Context, msg string, level Level) {
	data, _ := json.Marshal(LogData{
		Level:      level,
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
