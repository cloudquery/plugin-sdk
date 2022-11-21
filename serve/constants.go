package serve

import (
	"time"
)

const (
	// bufSize used for unit testing grpc server and client
	testBufSize       = 1024 * 1024
	flushTimeout      = 5 * time.Second
	maxReceiveMsgSize = 20 * 1024 * 1024 // 20MiB
)
