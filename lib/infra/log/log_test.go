package log_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ditointernet/go-dito/lib/infra/log"
)

func TestNewLogger(t *testing.T) {
	t.Run("should use LogLevel INFO when not specified", func(t *testing.T) {
		ctx := context.Background()

		logger := log.NewLogger(log.LoggerInput{})

		out := captureOutput(func() {
			logger.Debug(ctx, "random message")
		})

		if diff := cmp.Diff("", out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}

		out = captureOutput(func() {
			logger.Info(ctx, "random message")
		})

		if diff := cmp.Diff(`{"level":"INFO","message":"random message"}`, out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}
	})
}

func TestDebug(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       log.Level
		attrs       log.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message",
			expectedLog: `{"level":"DEBUG","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is INFO",
			ctx:         ctx,
			level:       log.LevelInfo,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       log.LevelWarning,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       log.LevelError,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       log.LevelCritical,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"DEBUG","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       log.LevelDebug,
			msg:         "random message",
			attrs:       log.LogAttributeSet{"attr1": true},
			expectedLog: `{"level":"DEBUG","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{Level: tc.level, Attributes: tc.attrs})

			out := captureOutput(func() {
				logger.Debug(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       log.Level
		attrs       log.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       log.LevelInfo,
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       log.LevelWarning,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       log.LevelError,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       log.LevelCritical,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"INFO","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       log.LevelDebug,
			msg:         "random message",
			attrs:       log.LogAttributeSet{"attr1": true},
			expectedLog: `{"level":"INFO","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{Level: tc.level, Attributes: tc.attrs})

			out := captureOutput(func() {
				logger.Info(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWarning(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       log.Level
		attrs       log.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       log.LevelInfo,
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       log.LevelWarning,
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       log.LevelError,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       log.LevelCritical,
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       log.LevelDebug,
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"WARNING","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       log.LevelDebug,
			msg:         "random message",
			attrs:       log.LogAttributeSet{"attr1": true},
			expectedLog: `{"level":"WARNING","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{Level: tc.level, Attributes: tc.attrs})

			out := captureOutput(func() {
				logger.Warning(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestError(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       log.Level
		attrs       log.LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       log.LevelDebug,
			err:         errors.New("random error"),
			expectedLog: `{"level":"ERROR","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       log.LevelInfo,
			err:         errors.New("random error"),
			expectedLog: `{"level":"ERROR","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       log.LevelWarning,
			err:         errors.New("random error"),
			expectedLog: `{"level":"ERROR","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       log.LevelError,
			err:         errors.New("random error"),
			expectedLog: `{"level":"ERROR","message":"random error"}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       log.LevelCritical,
			err:         errors.New("random error"),
			expectedLog: "",
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       log.LevelDebug,
			err:         errors.New("random error"),
			attrs:       log.LogAttributeSet{"attr1": true},
			expectedLog: `{"level":"ERROR","message":"random error","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{Level: tc.level, Attributes: tc.attrs})

			out := captureOutput(func() {
				logger.Error(tc.ctx, tc.err)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestCritical(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       log.Level
		attrs       log.LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       log.LevelDebug,
			err:         errors.New("random error"),
			expectedLog: `{"level":"CRITICAL","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       log.LevelInfo,
			err:         errors.New("random error"),
			expectedLog: `{"level":"CRITICAL","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       log.LevelWarning,
			err:         errors.New("random error"),
			expectedLog: `{"level":"CRITICAL","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       log.LevelError,
			err:         errors.New("random error"),
			expectedLog: `{"level":"CRITICAL","message":"random error"}`,
		},
		{
			desc:        "should log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       log.LevelCritical,
			err:         errors.New("random error"),
			expectedLog: `{"level":"CRITICAL","message":"random error"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       log.LevelDebug,
			err:         errors.New("random error"),
			attrs:       log.LogAttributeSet{"attr1": true},
			expectedLog: `{"level":"CRITICAL","message":"random error","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := log.NewLogger(log.LoggerInput{Level: tc.level, Attributes: tc.attrs})

			out := captureOutput(func() {
				logger.Critical(tc.ctx, tc.err)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func captureOutput(output func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return strings.TrimRight(string(out), "\n")
}
