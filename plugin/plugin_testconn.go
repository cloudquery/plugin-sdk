package plugin

import (
	"context"

	"github.com/rs/zerolog"
)

type TestConnFailureCode string

const (
	TestConnFailureCodeUnknown            TestConnFailureCode = "UNKNOWN"
	TestConnFailureCodeInvalidSpec        TestConnFailureCode = "INVALID_SPEC"
	TestConnFailureCodeInvalidCredentials TestConnFailureCode = "INVALID_CREDENTIALS"
)

type TestConnError struct {
	Code    TestConnFailureCode
	Message error
}

func NewTestConnError(code TestConnFailureCode, err error) *TestConnError {
	if code == "" {
		code = TestConnFailureCodeUnknown
	}
	return &TestConnError{
		Code:    code,
		Message: err,
	}
}

func (e *TestConnError) Error() string {
	return e.Message.Error()
}

func (e *TestConnError) Unwrap() error {
	return e.Message
}

func (e *TestConnError) Is(err error) bool {
	if err2, ok := err.(*TestConnError); ok {
		return e.Code == err2.Code
	}
	return false
}

type ConnectionTester func(ctx context.Context, logger zerolog.Logger, spec []byte) error

func (p *Plugin) TestConnection(ctx context.Context, logger zerolog.Logger, spec []byte) error {
	return p.testConnFn(ctx, logger, spec)
}

func UnimplementedTestConnectionFn(context.Context, zerolog.Logger, []byte) error {
	return ErrNotImplemented
}

var _ ConnectionTester = UnimplementedTestConnectionFn
