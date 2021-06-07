package brand

import (
	"strings"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/infra"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyBrandID is the key used to retrieve and save brand into the context
const ContextKeyBrandID string = "Brand"
const (
	// CodeTypeMissingBrand indicates that brand header is no present on the request
	CodeTypeMissingBrand errors.CodeType = "MISSING_BRAND"
	// CodeTypeMissingLoggerDependency indicates that logger dependency was not provided
	CodeTypeMissingLoggerDependency errors.CodeType = "MISSING_LOGGER_DEPENDENCY"
)

// AccountAuthenticator structure responsible for handling request authentication
type BrandFiller struct {
	logger infra.Logger
}

// BrandFiller creates a new instance of the Brand structure
func NewBrandFiller(logger infra.Logger) (BrandFiller, error) {
	if logger == nil {
		return BrandFiller{}, errors.New("missing logger dependency").WithKind(errors.KindInternal).WithCode(CodeTypeMissingLoggerDependency)
	}
	return BrandFiller{
		logger: logger,
	}, nil
}

// BrandFiller is the middleware responsible for retrieving brand id from the headers
//
func (ua BrandFiller) Fill(ctx *routing.Context) error {
	brandID := string(ctx.Request.Header.Peek(ContextKeyBrandID))
	brandID = strings.TrimSpace(brandID)
	if len(brandID) == 0 {
		err := errors.New("Brand is not present on request headers").WithKind(errors.KindUnauthorized).WithCode(CodeTypeMissingBrand)
		ua.logger.Error(ctx, err)
		return err
	}

	ctx.SetUserValue(ContextKeyBrandID, brandID)
	return nil
}
