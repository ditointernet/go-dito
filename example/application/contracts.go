package application

import (
	"context"

	"example/infra"
)

// Executable ...
type Executable interface {
	Run(context.Context) <-chan *infra.Error
}
