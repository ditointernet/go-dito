package brand

import (
	"errors"
	"net"
	"testing"

	ditoError "github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/mocks"
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

func newCtxWithHeaders(headers map[string]string) *routing.Context {
	ctx := newCtx()

	for key, value := range headers {
		ctx.Request.Header.Set(key, value)
	}

	return ctx
}

func TestBrandFillerFill(t *testing.T) {
	expectedBrand := "a-brand"
	var logger *mocks.MockLogger
	withMock := func(runner func(t *testing.T, m BrandFiller)) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			logger = mocks.NewMockLogger(ctrl)
			middleware, _ := NewBrandFiller(logger)

			runner(t, middleware)
		}
	}
	t.Run("return and error when there isnt a logger dependency injected on the constructor", func(t *testing.T) {
		_, err := NewBrandFiller(nil)

		var e ditoError.CustomError

		assert.True(t, errors.As(err, &e))
		assert.EqualError(t, e, "missing logger dependency")
	})
	t.Run("should include the brand in the context when brand is set in headers",
		withMock(func(t *testing.T, m BrandFiller) {
			ctx := newCtxWithHeaders(map[string]string{ContextKeyBrandID: expectedBrand})

			m.Fill(ctx)
			brand := ctx.Value(ContextKeyBrandID)

			assert.Equal(t, expectedBrand, brand)
		}))

	t.Run("should include space trimmed brand in the context when brand in the header has spaces in beginning or end",
		withMock(func(t *testing.T, m BrandFiller) {
			ctx := newCtxWithHeaders(map[string]string{ContextKeyBrandID: " a-brand    "})

			m.Fill(ctx)

			brand, ok := ctx.UserValue(ContextKeyBrandID).(string)

			assert.True(t, ok)
			assert.Equal(t, expectedBrand, brand)
		}))

	t.Run("should not include anything in the context when the brand in the headers is empty or has only spaces",
		withMock(func(t *testing.T, m BrandFiller) {
			ctx := newCtxWithHeaders(map[string]string{ContextKeyBrandID: " "})
			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			m.Fill(ctx)
			brand := ctx.UserValue(ContextKeyBrandID)

			assert.Nil(t, brand)
		}))

	t.Run("should not include anything in the context when there is no brand in the headers",
		withMock(func(t *testing.T, m BrandFiller) {
			ctx := newCtx()
			logger.EXPECT().Error(gomock.Any(), gomock.Any())
			m.Fill(ctx)
			brand := ctx.UserValue(ContextKeyBrandID)

			assert.Nil(t, brand)
		}))
}
