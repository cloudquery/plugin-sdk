package testlog

import (
	"io"
	"log"
	"testing"

	"github.com/hashicorp/go-hclog"
)

type testLogger struct {
	t     testing.TB
	level hclog.Level
}

// New create a new hclog adapter from testing.Log
// This is useful when running tests not to trash the output in tests
// and print only logs for tests that fail
func New(t testing.TB) hclog.Logger {
	return &testLogger{
		t:     t,
		level: hclog.Debug,
	}
}

func (l *testLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	l.t.Helper()
	switch level {
	case hclog.NoLevel:
		return
	case hclog.Trace:
		l.Trace(msg, args...)
	case hclog.Debug:
		l.Debug(msg, args...)
	case hclog.Info:
		l.Info(msg, args...)
	case hclog.Warn:
		l.Warn(msg, args...)
	case hclog.Error:
		l.Error(msg, args...)
	}
}

func (l *testLogger) Trace(msg string, args ...interface{}) {
	l.t.Helper()
	if l.level == hclog.Trace {
		l.t.Log(convertMsgArgToInterface("[TRACE] "+msg, args)...)
	}
}

func (l *testLogger) Debug(msg string, args ...interface{}) {
	l.t.Helper()
	if l.IsDebug() {
		l.t.Log(convertMsgArgToInterface("[DEBUG] "+msg, args)...)
	}
}

func (l *testLogger) Info(msg string, args ...interface{}) {
	l.t.Helper()
	if l.IsInfo() {
		l.t.Log(convertMsgArgToInterface("[INFO] "+msg, args)...)
	}
}

func (l *testLogger) Warn(msg string, args ...interface{}) {
	l.t.Helper()
	if l.IsWarn() {
		l.t.Log(convertMsgArgToInterface("[WARN] "+msg, args)...)
	}
}

func (l *testLogger) Error(msg string, args ...interface{}) {
	l.t.Helper()
	if l.IsError() {
		l.t.Log(convertMsgArgToInterface("[ERROR] "+msg, args)...)
	}
}

func (l *testLogger) IsTrace() bool {
	return l.level <= hclog.Trace
}

func (l *testLogger) IsDebug() bool {
	return l.level <= hclog.Debug
}

func (l *testLogger) IsInfo() bool {
	return l.level <= hclog.Info
}

func (l *testLogger) IsWarn() bool {
	return l.level <= hclog.Warn
}

func (l *testLogger) IsError() bool {
	return l.level <= hclog.Error
}

// ImpliedArgs returns With key/value pairs
func (*testLogger) ImpliedArgs() []interface{} {
	return nil
}

func (l *testLogger) With(args ...interface{}) hclog.Logger {
	return l
}

func (*testLogger) Name() string {
	return "testLogger"
}

func (l *testLogger) Named(name string) hclog.Logger {
	return l
}

func (l *testLogger) ResetNamed(name string) hclog.Logger {
	return l
}

func (l *testLogger) SetLevel(level hclog.Level) {
	l.level = level
}

func (*testLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return nil
}

func (*testLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}

func convertMsgArgToInterface(msg string, args ...interface{}) []interface{} {
	var res []interface{}
	res = append(res, msg)
	res = append(res, args...)
	return res
}
