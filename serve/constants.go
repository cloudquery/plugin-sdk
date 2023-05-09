package serve

import (
	"time"
)

const (
	// bufSize used for unit testing grpc server and client
	testBufSize  = 1024 * 1024
	flushTimeout = 5 * time.Second
	MaxMsgSize   = 100 * 1024 * 1024 // 100 MiB
)
