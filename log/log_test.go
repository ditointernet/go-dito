package log

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mockedTimmer() func() time.Time {
	return func() time.Time {
		return time.Date(2021, 0, 1, 12, 0, 0, 0, time.UTC)
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("should use LogLevel INFO when not specified", func(t *testing.T) {
		ctx := context.Background()

		logger := NewLogger(LoggerInput{})
		logger.now = mockedTimmer()

		out := captureOutput(func() {
			logger.Debug(ctx, "random message")
		})

		if diff := cmp.Diff("", out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}

		out = captureOutput(func() {
			logger.Info(ctx, "random message")
		})

		if diff := cmp.Diff(`{"time":"2020-12-01T12:00:00Z","severity":"INFO","message":"random message"}`, out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}
	})
}

func TestDebug(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"DEBUG","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"DEBUG","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       LogAttributeSet{"attr1": true},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"DEBUG","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerInput{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

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
		level       string
		attrs       LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"INFO","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"INFO","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"INFO","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       LogAttributeSet{"attr1": true},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"INFO","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerInput{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

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
		level       string
		attrs       LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"WARNING","message":"random message"}`,
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"WARNING","message":"random message with dynamic data 1"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       LogAttributeSet{"attr1": true},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"WARNING","message":"random message","attributes":{"attr1":"value1"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerInput{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

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
		level       string
		attrs       LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"ERROR","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"ERROR","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"ERROR","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"ERROR","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: "",
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       LogAttributeSet{"attr1": true},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"ERROR","message":"random error","attributes":{"attr1":"value1","code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerInput{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

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
		level       string
		attrs       LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       LogAttributeSet{"attr1": true},
			expectedLog: `{"time":"2020-12-01T12:00:00Z","severity":"CRITICAL","message":"random error","attributes":{"attr1":"value1","code":"UNKNOWN","kind":"UNEXPECTED"}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerInput{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

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
