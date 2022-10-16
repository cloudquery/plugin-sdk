package serve

import (
	"time"

	"google.golang.org/grpc/test/bufconn"
)

const (
	serveShort = `Start plugin server`
	// bufSize used for unit testing grpc server and client
	testBufSize  = 1024 * 1024
	flushTimeout = 5 * time.Second
)

// lis used for unit testing grpc server and client
var testSourceListener *bufconn.Listener
var testDestinationListener *bufconn.Listener
