package opa

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/discovery"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"

	"github.com/ditointernet/go-dito/lib/errors"
)

// AuthorizationResult ...
type AuthorizationResult []map[string]interface{}

// Client is an OPA authorization client
type Client struct {
	manager *plugins.Manager
}

// NewClient creates a new Client object
func NewClient(opaBundleBaseURL string, opaPollingMinDelay, opaPollingMaxDelay int) (Client, error) {
	manager, err := plugins.New(buildOPAConfigFile(opaBundleBaseURL, opaPollingMinDelay, opaPollingMaxDelay), uuid.New().String(), inmem.New())
	if err != nil {
		return Client{}, errors.New(err.Error()).WithKind(errors.KindInvalidInput).WithCode(ERR_CODE_BUILD_MANAGER_FAILURE)
	}

	d, err := discovery.New(manager)
	if err != nil {
		return Client{}, errors.New(err.Error()).WithCode(ERR_CODE_BUILD_DISCOVERY_FAILURE)
	}

	manager.Register("discovery", d)
	if err = manager.Start(context.Background()); err != nil {
		return Client{}, errors.New(err.Error()).WithCode(ERR_CODE_START_MANAGER_FAILURE)
	}

	return Client{manager: manager}, nil
}

// MustNewClient creates a new Client object.
// It panics if any error is found.
func MustNewClient(opaBundleBaseURL string, opaPollingMinDelay, opaPollingMaxDelay int) Client {
	cli, err := NewClient(opaBundleBaseURL, opaPollingMinDelay, opaPollingMaxDelay)
	if err != nil {
		panic(err)
	}

	return cli
}

// DecideIfAllowed indicates if the given user is allowed to perform the given action
func (c Client) DecideIfAllowed(ctx context.Context, regoQuery, method, path, brandID, userID string) (bool, error) {
	var decision bool

	err := storage.Txn(ctx, c.manager.Store, storage.TransactionParams{}, func(txn storage.Transaction) error {
		var err error

		input := map[string]interface{}{
			"method":   method,
			"path":     path,
			"brand_id": brandID,
			"user_id":  userID,
		}

		q, err := rego.New(
			rego.Query(regoQuery),
			rego.Input(input),
			rego.Compiler(c.manager.GetCompiler()),
			rego.Store(c.manager.Store),
			rego.Transaction(txn),
		).PrepareForEval(ctx)
		if err != nil {
			return errors.New(err.Error()).WithCode(ERR_CODE_BUILD_REGO_OBJECT_FAILURE)
		}

		rs, err := q.Eval(ctx)
		if err != nil {
			return errors.New(err.Error()).WithCode(ERR_CODE_EVAL_REGO_FAILURE)
		}
		if len(rs) == 0 {
			return ErrNoDecision
		}

		var ok bool
		decision, ok = rs[0].Expressions[0].Value.(bool)
		if !ok || len(rs) > 1 {
			return ErrNonBooleanDecision
		}

		return nil
	})

	return decision, err
}

// ExecuteQuery ...
func (c Client) ExecuteQuery(ctx context.Context, query string, input map[string]interface{}) (AuthorizationResult, error) {
	var regoResult rego.ResultSet
	var queryResult AuthorizationResult

	err := storage.Txn(ctx, c.manager.Store, storage.TransactionParams{}, func(txn storage.Transaction) error {
		q, err := rego.New(
			rego.Query(query),
			rego.Input(input),
			rego.Compiler(c.manager.GetCompiler()),
			rego.Store(c.manager.Store),
			rego.Transaction(txn),
		).PrepareForEval(ctx)
		if err != nil {
			return errors.New(err.Error()).WithCode(ERR_CODE_BUILD_REGO_OBJECT_FAILURE)
		}

		if regoResult, err = q.Eval(ctx); err != nil {
			return errors.New(err.Error()).WithCode(ERR_CODE_EVAL_REGO_FAILURE)
		}

		return nil
	})

	for _, r := range regoResult {
		queryResult = append(queryResult, r.Bindings)
	}
	return queryResult, err
}

func buildOPAConfigFile(opaBundleBaseURL string, opaPollingMinDelay, opaPollingMaxDelay int) []byte {
	configFile := fmt.Sprintf(`---
services:
  opaapi:
    url: %s
    credentials:
      s3_signing:
        environment_credentials: {}
bundles:
  authz:
    service: opaapi
    resource: bundle.tar.gz
    polling:
      min_delay_seconds: %d
      max_delay_seconds: %d
`,
		opaBundleBaseURL, opaPollingMinDelay, opaPollingMaxDelay)

	return []byte(configFile)
}
