package plugin

import (
	"context"
)

type TestConnFailureCode string

const (
	TestConnFailureCodeUnknown            TestConnFailureCode = "UNKNOWN"
	TestConnFailureCodeUnimplemented      TestConnFailureCode = "UNIMPLEMENTED"
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

var ErrTestConnUnimplemented = &TestConnError{
	Code: TestConnFailureCodeUnimplemented,
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

type ConnectionTester func(ctx context.Context, spec []byte) *TestConnError

func (p *Plugin) TestConnection(ctx context.Context, spec []byte) *TestConnError {
	return p.testConnFn(ctx, spec)
}

func UnimplementedTestConnectionFn(context.Context, []byte) *TestConnError {
	return ErrTestConnUnimplemented
}

var _ ConnectionTester = UnimplementedTestConnectionFn
