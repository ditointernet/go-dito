package authorization

import (
	"context"
	"fmt"

	"time"

	"github.com/ditointernet/go-dito/lib/errors"
	"github.com/ditointernet/go-dito/lib/http/middleware/authentication"
	"github.com/ditointernet/go-dito/lib/http/middleware/brand"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// ContextKeyAllowedStores is the context key that get and sets all accounts allowed stores
const ContextKeyAllowedStores string = "allowed-stores"

// ResourseFilter defines a type that represents a resource filter as an integer
type ResourseFilter int

const (
	// CodeTypeMissingLoggerDependency indicates that the logger dependency parameter was not injected
	CodeTypeMissingLoggerDependency errors.CodeType = "MISSING_LOGGER_DEPENDENCY"
	// CodeTypeMissingAuthClientDependency indicates that the authclient dependency parameter was not injected
	CodeTypeMissingAuthClientDependency errors.CodeType = "MISSING_AUTH_CLIENT_DEPENDENCY"
	// CodeTypeMissingResourceNameDependency indicates that the resourceName dependency parameter was not injected
	CodeTypeMissingResourceNameDependency errors.CodeType = "MISSING_RESOURCE_NAME_DEPENDENCY"
	// CodeTypeMissingAccountId indicates that the accountId was not present on the context
	CodeTypeMissingAccountId errors.CodeType = "MISSING_ACCOUNT_ID"
	// CodeTypeMissingBrandId indicates that the brandId was not present on the context
	CodeTypeMissingBrandId errors.CodeType = "MISSING_BRAND_ID"
	// CodeTypeErrorExecutingAuthorizationQuery indicates that it was not possible to execute authorization client query
	CodeTypeErrorExecutingAuthorizationQuery errors.CodeType = "CANT_EXECUTE_AUTH_QUERY"
	// CodeTypeAccessDenied indicates that authorizatior client denied account access
	CodeTypeAccessDenied errors.CodeType = "ACCESS_DENIED"
)

const (
	// StoreFilter means a numeric representation of the stores filter
	StoreFilter ResourseFilter = iota
)

// String returns the equivalent string of the Resource filters integers
func (s ResourseFilter) String() string {
	return [...]string{
		"stores",
	}[s]
}

// AccountAuthorizator is the struct responsible for create account authorizarion
type AccountAuthorizator struct {
	logger              logger
	authorizatorClient  authorizatorClient
	authorizatorTimeout time.Duration
	Now                 func() time.Time
	resourceName        string
	ResourseFilters     []ResourseFilter
}

// NewAccountAuthorizator constructs a new account authorization middleware
func NewAccountAuthorizator(
	logger logger,
	authClient authorizatorClient,
	authorizatorTimeout time.Duration,
	resourceName string,
	ResourseFilters []ResourseFilter,
) (AccountAuthorizator, error) {
	if resourceName == "" {
		err := errors.New("missing resource name").WithKind(errors.KindInternal).WithCode(CodeTypeMissingResourceNameDependency)
		return AccountAuthorizator{}, err
	}

	if authClient == nil {
		err := errors.New("missing auth client").WithKind(errors.KindInternal).WithCode(CodeTypeMissingAuthClientDependency)
		return AccountAuthorizator{}, err
	}

	if logger == nil {
		err := errors.New("missing logger").WithKind(errors.KindInternal).WithCode(CodeTypeMissingLoggerDependency)
		return AccountAuthorizator{}, err
	}

	return AccountAuthorizator{
		logger:              logger,
		authorizatorClient:  authClient,
		authorizatorTimeout: authorizatorTimeout,
		resourceName:        resourceName,
		ResourseFilters:     ResourseFilters,
		// Just for allowing mock time
		Now: time.Now,
	}, nil
}

// Authorize is the middleware responsible for calling the auth client and check if user is authorized to make the current request
func (a AccountAuthorizator) Authorize(ctx *routing.Context) error {
	accountID := ctx.Value(authentication.ContextKeyAccountID)
	if accountID == nil {
		err := errors.New("missing account id").WithCode(CodeTypeMissingAccountId)
		a.logger.Error(ctx, err)
		return err
	}
	brandID := ctx.Value(brand.ContextKeyBrandID)
	if brandID == nil {
		err := errors.New("missing brand id").WithCode(CodeTypeMissingBrandId)
		a.logger.Error(ctx, err)
		return err
	}
	c, cancel := context.WithDeadline(ctx, a.Now().Add(a.authorizatorTimeout))
	defer cancel()

	query := fmt.Sprintf(`allow := data.authz.%s_allow`, a.resourceName)

	resourceInput := map[string]interface{}{
		"method":   string(ctx.Method()),
		"path":     string(ctx.Path()),
		"brand_id": brandID,
		"user_id":  accountID,
	}
	if len(a.ResourseFilters) > 0 {
		query = query + ` ; filter := data.authz.filter_values`
		resourceInput["filter_type"] = a.ResourseFilters[0].String()
	}
	result, err := a.authorizatorClient.ExecuteQuery(c, query, resourceInput)
	if err != nil {
		err := errors.New("error on executing opa client query, got : %s", err).WithKind(errors.KindInternal).WithCode(CodeTypeErrorExecutingAuthorizationQuery)
		a.logger.Error(ctx, err)
		return err
	}
	if len(result) == 0 {
		err := errors.New("error on executing authorizator client query, got undefined result: %s", err).WithKind(errors.KindInternal).WithCode(CodeTypeErrorExecutingAuthorizationQuery)
		a.logger.Error(ctx, err)
		return err
	}

	allowed, ok := result[0]["allow"].(bool)
	if !ok {
		err := errors.New("error on executing authorizator client query, allow condition not found").WithKind(errors.KindInternal).WithCode(CodeTypeErrorExecutingAuthorizationQuery)
		a.logger.Error(ctx, err)
		return err
	}

	if !allowed {
		err := errors.New("Authorization decision - accountID: %s with brandID %s access was denied", accountID, brandID).WithKind(errors.KindUnauthorized).WithCode(CodeTypeAccessDenied)
		a.logger.Error(ctx, err)
		return err
	}

	if len(a.ResourseFilters) > 0 {
		filterValues, _ := result[0]["filter"].([]interface{})
		var allowedStores []string

		for _, f := range filterValues {
			store, ok := f.(string)
			if ok {
				allowedStores = append(allowedStores, store)
			}
		}

		ctx.SetUserValue(ContextKeyAllowedStores, allowedStores)
	}

	a.logger.Debug(ctx, "Authorization decision - accountID: %s with brandID %s access was granted", accountID, brandID)
	return nil
}
