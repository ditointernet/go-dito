package brand

import (
	"strings"

	"github.com/ditointernet/go-dito/lib/errors"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyBrandID is the key used to retrieve and save brand into the context
const ContextKeyBrandID string = "Brand"

// AccountAuthenticator structure responsible for handling request authentication
type BrandFiller struct {
	logger logger
}

// BrandFiller creates a new instance of the Brand structure
func NewBrandFiller(logger logger) (BrandFiller, error) {
	if logger == nil {
		return BrandFiller{}, errors.New("missing logger dependency").WithKind(errors.KindInternal)
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
		err := errors.New("Brand is not present on request headers").WithKind(errors.KindUnauthorized)
		ua.logger.Error(ctx, err)
		return err
	}

	ctx.SetUserValue(ContextKeyBrandID, brandID)
	return nil
}
