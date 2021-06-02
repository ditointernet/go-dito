package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
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
		l.printError(ctx, err, LevelError)
	}
}

// Critical logs critical data
func (l Logger) Critical(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.printError(ctx, err, LevelCritical)
	}
}

type logData struct {
	TraceID    string                       `json:"trace_id,omitempty"`
	Timestamp  string                       `json:"timestamp"`
	Level      string                       `json:"level"`
	Message    string                       `json:"message"`
	Attributes map[LogAttribute]interface{} `json:"attributes,omitempty"`
}

func (l Logger) print(ctx context.Context, msg string, level Level) {
	span := trace.SpanFromContext(ctx)

	attrs := l.extractLogAttributesFromContext(ctx)

	span.AddEvent("log", trace.WithAttributes(buildOtelAttributes(attrs, "log")...))

	data, _ := json.Marshal(logData{
		TraceID:    span.SpanContext().TraceID.String(),
		Timestamp:  l.now().Format(time.RFC3339),
		Level:      level.String(),
		Message:    msg,
		Attributes: attrs,
	})

	fmt.Println(string(data))
}

func (l Logger) printError(ctx context.Context, err error, level Level) {
	span := trace.SpanFromContext(ctx)

	attrs := l.extractLogAttributesFromContext(ctx)

	span.RecordError(err, trace.WithAttributes(buildOtelAttributes(attrs, "exception")...))

	data, _ := json.Marshal(logData{
		TraceID:    span.SpanContext().TraceID.String(),
		Timestamp:  l.now().Format(time.RFC3339),
		Level:      level.String(),
		Message:    err.Error(),
		Attributes: attrs,
	})

	fmt.Println(string(data))
}

func buildOtelAttributes(attrs map[LogAttribute]interface{}, prefix string) []label.KeyValue {
	eAttrs := []label.KeyValue{}
	for k, v := range attrs {
		eAttrs = append(eAttrs, label.String(fmt.Sprintf("%s.%s", prefix, k), v.(string)))
	}

	return eAttrs
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
