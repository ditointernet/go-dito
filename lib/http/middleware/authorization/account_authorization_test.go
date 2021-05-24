package authorization

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	ditoError "github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/middleware/authentication"
	"github.com/ditointernet/go-dito/lib/http/middleware/authorization/mocks"
	"github.com/ditointernet/go-dito/lib/http/middleware/brand"
	opaLib "github.com/ditointernet/go-dito/lib/opa"
	"github.com/golang/mock/gomock"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func newCtx() *routing.Context {
	req := fasthttp.AcquireRequest()
	reqCtx := &fasthttp.RequestCtx{}
	reqCtx.Init(req, &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}, nil)
	return &routing.Context{
		RequestCtx: reqCtx,
	}
}

func newCtxWithUserValues(userValues map[string]interface{}) *routing.Context {
	ctx := newCtx()

	for key, value := range userValues {
		ctx.SetUserValue(key, value)
	}

	return ctx
}

func TestAuthorize(t *testing.T) {
	var logger *mocks.Mocklogger
	var opa *mocks.MockauthorizatorClient
	timeout := 100 * time.Millisecond

	withMock := func(runner func(t *testing.T, m AccountAuthorizator)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger = mocks.NewMocklogger(ctrl)
			opa = mocks.NewMockauthorizatorClient(ctrl)
			middleware, _ := NewAccountAuthorizator(logger, opa, timeout, "some-client", []ResourseFilter{})

			runner(t, middleware)
		}
	}
	t.Run("should not create the authorizator instance when there isn't a resource name",
		func(t *testing.T) {
			_, err := NewAccountAuthorizator(logger, opa, timeout, "", []ResourseFilter{})

			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing resource name")
		})
	t.Run("should not create the authorizator instance when there isn't a logger",
		func(t *testing.T) {
			_, err := NewAccountAuthorizator(nil, opa, timeout, "some-client", []ResourseFilter{})

			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing logger")
		})
	t.Run("should not create the authorizator instance when there isn't a authCLient",
		func(t *testing.T) {
			_, err := NewAccountAuthorizator(logger, nil, timeout, "some-client", []ResourseFilter{})

			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing auth client")
		})

	t.Run("should not authorize when there is no user id on context",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{brand.ContextKeyBrandID: "any-brand2"})

			logger.EXPECT().Error(gomock.Any(), gomock.Any())

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing user id")
		}))

	t.Run("should not authorize when there is no brand id on headers",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{authentication.ContextKeyAccountID: "123456"})

			logger.EXPECT().Error(gomock.Any(), gomock.Any())

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing brand id")
		}))

	t.Run("should return error when opa client returns an error",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error"))

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing opa client query, got : any error")
		}))

	t.Run("should return error when opa client returns an empty result",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(nil, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing authorizator client query, got undefined result: %!s(<nil>)")
		}))

	t.Run("should return error when allow response is not found",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{}}, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing authorizator client query, allow condition not found")
		}))

	t.Run("should not authorize unauthorized users",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand4",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{"allow": false}}, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "Authorization decision - accountID: 123456 with brandID any-brand4 access was denied")
		}))

	t.Run("should authorize user when the user exists in the bundle data",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand1",
			})
			timeout := 100 * time.Millisecond
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()
			logger.EXPECT().Info(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{"allow": true}}, nil)
			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)

			assert.Nil(t, err)
		}))
}

func TestAuthorize_WithFIlters(t *testing.T) {
	var logger *mocks.Mocklogger
	var opa *mocks.MockauthorizatorClient
	timeout := 100 * time.Millisecond

	withMock := func(runner func(t *testing.T, m AccountAuthorizator)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger = mocks.NewMocklogger(ctrl)
			opa = mocks.NewMockauthorizatorClient(ctrl)
			middleware, _ := NewAccountAuthorizator(logger, opa, timeout, "some-client", []ResourseFilter{StoreFilter})

			runner(t, middleware)
		}
	}

	t.Run("should not authorize when there is no user id on context",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{brand.ContextKeyBrandID: "any-brand2"})

			logger.EXPECT().Error(gomock.Any(), gomock.Any())

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing user id")
		}))

	t.Run("should not authorize when there is no brand id on headers",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{authentication.ContextKeyAccountID: "123456"})

			logger.EXPECT().Error(gomock.Any(), gomock.Any())

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "missing brand id")
		}))

	t.Run("should return error when opa client returns an error",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error"))

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing opa client query, got : any error")
		}))

	t.Run("should return error when opa client returns an empty result",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(nil, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing authorizator client query, got undefined result: %!s(<nil>)")
		}))

	t.Run("should return error when allow response is not found",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand3",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{}}, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			assert.EqualError(t, err, "error on executing authorizator client query, allow condition not found")
		}))

	t.Run("should not authorize unauthorized users",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand4",
			})
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()

			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{"allow": false}}, nil)

			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)
			if err == nil {
				t.Fatal("expected error was not found")
			}

			var e ditoError.CustomError
			assert.True(t, errors.As(err, &e))
			assert.EqualError(t, e, "Authorization decision - accountID: 123456 with brandID any-brand4 access was denied")
		}))

	t.Run("should authorize user when the user exists in the bundle data",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand1",
			})
			timeout := 100 * time.Millisecond
			now, _ := time.Parse(time.RFC3339, "2020-01-30T03:00:00Z")
			expectedCtx, cancel := context.WithDeadline(ctx, now.Add(timeout))
			defer cancel()
			logger.EXPECT().Info(gomock.Any(), gomock.Any())
			opa.EXPECT().
				ExecuteQuery(expectedCtx, gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{{"allow": true}}, nil)
			m.Now = func() time.Time { return now }

			err := m.Authorize(ctx)

			assert.Nil(t, err)
		}))

	t.Run("should set filter values when they are found",
		withMock(func(t *testing.T, m AccountAuthorizator) {
			logger.EXPECT().Info(gomock.Any(), gomock.Any())

			ctx := newCtxWithUserValues(map[string]interface{}{
				authentication.ContextKeyAccountID: "123456",
				brand.ContextKeyBrandID:            "any-brand1",
			})

			opa.EXPECT().
				ExecuteQuery(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(opaLib.AuthorizationResult{
					{"allow": true, "filter": []interface{}{"fil1", "fil2"}},
				}, nil)

			m.Authorize(ctx)

			expected := []string{"fil1", "fil2"}
			received, _ := ctx.UserValue(ContextKeyAllowedStores).([]string)
			assert.ElementsMatch(t, expected, received)
		}))
}
