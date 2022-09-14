package steps_test

import (
	"testing"

	"github.com/ditointernet/go-dito/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ErrMock is a mocked error. Should only be used for testing purposes.
var ErrMock = errors.New("mocked error")

func TestSteps(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Steps Suite")
}

func fillIntegerChannel(ch chan any, numItems int) {
	go func() {
		for i := 0; i < numItems; i++ {
			ch <- i
		}
	}()
}
