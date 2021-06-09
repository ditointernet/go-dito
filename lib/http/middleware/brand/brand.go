package brand

import (
	"strings"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/infra"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyBrandID is the key used to retrieve and save brand into the context
const ContextKeyBrandID string = "brand_id"

const (
	// CodeTypeMissingBrand indicates that brand header is no present on the request
	CodeTypeMissingBrand errors.CodeType = "MISSING_BRAND"
)

// BrandFiller structure responsible for injecting the brand into HTTP request context.
type BrandFiller struct {
	logger infra.Logger
}

// NewBrandFiller creates a new instance of the Brand structure
func NewBrandFiller(logger infra.Logger) (BrandFiller, error) {
	if logger == nil {
		return BrandFiller{}, errors.NewMissingRequiredDependency("logger")
	}

	return BrandFiller{logger: logger}, nil
}

// MustNewBrandFiller creates a new instance of the Brand structure.
// It panics if any error is found.
func MustNewBrandFiller(logger logger) BrandFiller {
	mid, err := NewBrandFiller(logger)
	if err != nil {
		panic(err)
	}

	return mid
}

// Fill is the middleware responsible for retrieving brand id from the headers.
func (ua BrandFiller) Fill(ctx *routing.Context) error {
	brandID := string(ctx.Request.Header.Peek("brand"))
	brandID = strings.TrimSpace(brandID)
	if len(brandID) == 0 {
		err := errors.New("brand is not present on request headers").WithKind(errors.KindUnauthorized).WithCode(CodeTypeMissingBrand)
		ua.logger.Error(ctx, err)
		return err
	}

	ctx.SetUserValue(ContextKeyBrandID, brandID)
	return nil
}
