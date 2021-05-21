package authorization

import (
	"context"
	"fmt"

	"time"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/middleware/authentication"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyAllowedStores is the context key that get and sets all accounts allowed stores
const ContextKeyAllowedStores string = "allowed-stores"

// AccountAuthorizator is the struct responsible for create account authorizarion
type AccountAuthorizator struct {
	logger              Logger
	authorizatorClient  AuthorizatorClient
	authorizatorTimeout time.Duration
	Now                 func() time.Time
	regoClient          string
}

// NewAccountAuthorizator constructs a new account authorization middleware
func NewAccountAuthorizator(
	logger Logger,
	authClient AuthorizatorClient,
	authorizatorTimeout time.Duration,
	regoClient string,
) (AccountAuthorizator, error) {
	return AccountAuthorizator{
		logger:              logger,
		authorizatorClient:  authClient,
		authorizatorTimeout: authorizatorTimeout,
		regoClient:          regoClient,
		// Just for allowing mock time
		Now: time.Now,
	}, nil
}

// Authorize is the middleware responsible for calling the auth client and check if user is authorized to make the current request
func (a AccountAuthorizator) Authorize(ctx *routing.Context) error {

	accountID := ctx.Value(authentication.ContextKeyAccountID)
	if accountID == nil {
		err := errors.New("missing user id").WithKind(errors.KindInternal)
		a.logger.Error(ctx, err)
		return err
	}
	// todo get brand id from package brand id
	brandID := ctx.Value("brand-id")
	if brandID == nil {
		err := errors.New("missing brand id").WithKind(errors.KindInternal)
		a.logger.Error(ctx, err)
		return err
	}

	c, cancel := context.WithDeadline(ctx, a.Now().Add(a.authorizatorTimeout))
	defer cancel()

	query := fmt.Sprintf(`allow := data.authz.%s_allow ; filter := data.authz.filter_values`, a.regoClient)
	result, err := a.authorizatorClient.ExecuteQuery(c, query, map[string]interface{}{
		"method":      string(ctx.Method()),
		"path":        string(ctx.Path()),
		"brand_id":    brandID,
		"user_id":     accountID,
		"filter_type": "stores",
	})
	if err != nil {
		err := errors.New("error on executing opa client query, got : %s", err).WithKind(errors.KindInternal)
		a.logger.Error(ctx, err)
		return err
	}
	if len(result) == 0 {
		err := errors.New("error on executing authorizator client query, got undefined result: %s", err).WithKind(errors.KindInternal)
		a.logger.Error(ctx, err)
		return err
	}

	allowed, ok := result[0]["allow"].(bool)
	if !ok {
		err := errors.New("error on executing authorizator client query, allow condition not found").WithKind(errors.KindInternal)
		a.logger.Error(ctx, err)
		return err
	}

	if !allowed {
		err := errors.New("Authorization decision - accountID: %s with brandID %s access was denied", accountID, brandID).WithKind(errors.KindUnauthorized)
		a.logger.Error(ctx, err)
		return err
	}

	filterValues, _ := result[0]["filter"].([]interface{})
	var allowedStores []string

	for _, f := range filterValues {
		store, ok := f.(string)
		if ok {
			allowedStores = append(allowedStores, store)
		}
	}

	ctx.SetUserValue(ContextKeyAllowedStores, allowedStores)

	a.logger.Info(ctx, "Authorization decision - accountID: %s with brandID %s access granted")
	return nil
}
