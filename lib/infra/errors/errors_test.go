package errors_test

import (
	e "errors"
	"testing"

	"github.com/ditointernet/go-dito/lib/infra/errors"
)

// TestNew tests a creation of a new error and its mutations
func TestNew(t *testing.T) {
	t.Run("should produce and error with the given message", func(t *testing.T) {
		msg := "mocked message"
		if err := errors.New(msg); err.Error() != msg {
			t.Errorf("Wrong error message received. Expected '%s', got '%s'", msg, err.Error())
		}
	})

	t.Run("should produce and error with an dynamic message", func(t *testing.T) {
		expectedMsg := "mocked message with dynamic value 1000"
		if err := errors.New("mocked message with dynamic value %d", 1000); err.Error() != expectedMsg {
			t.Errorf("Wrong error message received. Expected '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("should produce and error with the default kind and code when not informed", func(t *testing.T) {
		err := errors.New("mocked message")

		if errors.Kind(err) != errors.KindUnexpected {
			t.Errorf("expected '%s', got '%s'", errors.KindUnexpected, errors.Kind(err))
		}

		if errors.Code(err) != errors.CodeUnknown {
			t.Errorf("expected '%s', got '%s'", errors.CodeUnknown, errors.Code(err))
		}
	})

	t.Run("should produce and error with a non-defaul Kind filled", func(t *testing.T) {
		err := errors.New("mocked message").WithKind(errors.KindInternal)

		if errors.Kind(err) != errors.KindInternal {
			t.Errorf("expected '%s', got '%s'", errors.KindInternal, errors.Kind(err))
		}
	})

	t.Run("should override the Kind if WithKind is called more than once", func(t *testing.T) {
		err := errors.New("mocked message").WithKind(errors.KindInternal).WithKind(errors.KindNotFound)

		if errors.Kind(err) != errors.KindNotFound {
			t.Errorf("expected '%s', got '%s'", errors.KindNotFound, errors.Kind(err))
		}
	})

	t.Run("should produce and error with a non-defaul Code filled", func(t *testing.T) {
		err := errors.New("mocked message").WithCode("MOCKED_CODE")

		if errors.Code(err) != "MOCKED_CODE" {
			t.Errorf("expected 'MOCKED_CODE', got %s", errors.Code(err))
		}
	})

	t.Run("should override the Code if WithCode is called more than once", func(t *testing.T) {
		err := errors.New("mocked message").WithCode("MOCKED_CODE").WithCode("MOCKED_CODE_2")

		if errors.Code(err) != "MOCKED_CODE_2" {
			t.Errorf("expected 'MOCKED_CODE_2', got %s", errors.Code(err))
		}
	})
}

func TestKind(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedKind errors.KindType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedKind: errors.KindUnexpected,
		},
		{
			name:         "custom error with default kind",
			err:          errors.New("some message"),
			expectedKind: errors.KindUnexpected,
		},
		{
			name:         "custom error with non-default kind",
			err:          errors.New("some message").WithKind("some kind"),
			expectedKind: "some kind",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if kind := errors.Kind(tc.err); kind != tc.expectedKind {
				t.Errorf("Expected kind to be '%s': received '%s'", tc.expectedKind, kind)
			}
		})
	}
}

func TestCode(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedCode errors.CodeType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedCode: errors.CodeUnknown,
		},
		{
			name:         "custom error with default code",
			err:          errors.New("some message"),
			expectedCode: errors.CodeUnknown,
		},
		{
			name:         "custom error with non-default code",
			err:          errors.New("some message").WithCode("some code"),
			expectedCode: "some code",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if code := errors.Code(tc.err); code != tc.expectedCode {
				t.Errorf("Expected code to be: '%s', received '%s'", tc.expectedCode, code)
			}
		})
	}
}
