package plugin

import (
	"context"
	"fmt"
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

// ConnectionTester is an interface that can be implemented by a plugin client to enable explicit connection testing.
type ConnectionTester interface {
	TestConnection(context.Context, []byte) *TestConnError
}

func (p *Plugin) TestConnection(ctx context.Context, spec []byte) *TestConnError {
	if !p.mu.TryLock() {
		return &TestConnError{
			Code:    TestConnFailureCodeUnknown,
			Message: fmt.Errorf("plugin already in use"),
		}
	}
	defer p.mu.Unlock()

	if v, ok := p.client.(ConnectionTester); ok {
		return v.TestConnection(ctx, spec)
	}

	return ErrTestConnUnimplemented
}
