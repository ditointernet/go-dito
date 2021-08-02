package authentication

import "github.com/ditointernet/go-dito/lib/errors"

var (
	// ErrMissingOrInvalidAuthenticationToken ...
	ErrMissingOrInvalidAuthenticationToken = errors.New("missing or invalid authentication token").
						WithKind(errors.KindUnauthenticated).
						WithCode("MISSING_OR_INVALID_AUTHENTICATION_TOKEN")

	// ErrRenewCertificates ...
	ErrRenewCertificates = errors.New("error on renewing the certificates").
				WithKind(errors.KindInternal).
				WithCode("FAILED_RENEWING_JKWS_CERTIFICATES")

	// ErrInvalidJWTToken ...
	ErrInvalidJWTToken = errors.New("invalid JWT token").
				WithKind(errors.KindUnauthenticated).
				WithCode("INVALID_JWT_TOKEN")
)
