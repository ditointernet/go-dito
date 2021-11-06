package authentication

import "context"

type JWKSClient interface {
	GetCerts(ctx context.Context) error
	RenewCerts(ctx context.Context) error
	Certs() map[string]string
}
