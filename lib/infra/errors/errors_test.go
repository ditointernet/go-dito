package errors

import (
	e "errors"
	"fmt"
	"testing"
)

// TestNew tests a creation of a new error with only the kind parameter
func TestNewType(t *testing.T) {
	var customError CustomError

	t.Run("Should return an error with type Custom Error", func(t *testing.T) {

		err := New("some message", "some kind", "some code")

		if !e.As(err, &customError) {
			t.Errorf("Expected error as a type Custom Error, got: '%T'", err)
		}
	})

}

func TestNewValues(t *testing.T) {
	err := New("some message", "some kind", "some code")
	t.Run("should have same kind of the new error", func(t *testing.T) {
		var expectedKind KindType
		expectedKind = "some kind"
		if expectedKind != err.kind {
			t.Errorf("expected kind to be: '%s' got: '%s'", expectedKind, err.kind)
		}
	})
	t.Run("should have same message of the new error", func(t *testing.T) {
		var expectedMessage string
		expectedMessage = "some message"
		if expectedMessage != err.message {
			t.Errorf("expected message to be: '%s' got: '%s'", expectedMessage, err.kind)
		}
	})
	t.Run("should have same kind of the new error", func(t *testing.T) {
		var expectedCode CodeType
		expectedCode = "some code"
		if expectedCode != err.code {
			t.Errorf("expected code to be: '%s' got: '%s'", expectedCode, err.code)
		}
	})
}

func TestError(t *testing.T) {
	err := New("some message", "some kind", "some code")
	t.Run("should have same kind of the new error", func(t *testing.T) {
		expectedMessage := "some message"
		if expectedMessage != err.Error() {
			t.Errorf("expected message to be: '%s' got: '%s'", expectedMessage, err.kind)
		}
	})
}
func TestKind(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedKind KindType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedKind: DefaultKind,
		},
		{
			name:         "custom error",
			err:          New("some message", "some kind", ""),
			expectedKind: "some kind",
		},
		{
			name:         "empty kind",
			err:          New("some message", "", ""),
			expectedKind: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			kind := Kind(tc.err)
			fmt.Println(kind)
			if kind != tc.expectedKind {
				t.Errorf("Expected kind to be '%s': received '%s'", tc.expectedKind, kind)
			}
		})
	}
}

func TestCode(t *testing.T) {
	tt := []struct {
		name         string
		err          error
		expectedCode CodeType
	}{
		{
			name:         "go native error",
			err:          e.New("new error"),
			expectedCode: DefaultCode,
		},
		{
			name:         "custom error",
			err:          New("some message", "some kind", "some code"),
			expectedCode: "some code",
		},
		{
			name:         "empty code",
			err:          New("some message", "", ""),
			expectedCode: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			code := Code(tc.err)
			if code != tc.expectedCode {
				t.Errorf("Expected code to be: '%s', received '%s'", tc.expectedCode, code)
			}
		})
	}
}
