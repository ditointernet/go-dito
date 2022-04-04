package opa

import (
	"github.com/ditointernet/go-dito/lib/errors"
)

const (
	ERR_CODE_BUILD_MANAGER_FAILURE     = "OPA_BUILD_MANAGER_FAILURE"
	ERR_CODE_BUILD_DISCOVERY_FAILURE   = "OPA_BUILD_DISCOVERY_FAILURE"
	ERR_CODE_START_MANAGER_FAILURE     = "OPA_START_MANAGER_FAILURE"
	ERR_CODE_BUILD_REGO_OBJECT_FAILURE = "OPA_BUILD_REGO_FAILURE"
	ERR_CODE_EVAL_REGO_FAILURE         = "OPA_EVAL_REGO_FAILURE"
	ERR_CODE_NO_DECISION               = "OPA_NO_DECISION"
	ERR_CODE_NON_BOOLEAN_DECISION      = "OPA_NON_BOOLEAN_DECISION"
)

var (
	ErrNoDecision                            = errors.New("Undefined decision").WithCode(ERR_CODE_NO_DECISION)
	ErrNonBooleanDecision errors.CustomError = errors.New("Non-boolean decision").WithCode(ERR_CODE_NON_BOOLEAN_DECISION)
)
