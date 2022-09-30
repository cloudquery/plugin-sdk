package serve

import (
	"time"

	"github.com/cloudquery/plugin-sdk/plugins"
	"google.golang.org/grpc/test/bufconn"
)

type Options struct {
	// Required: Source or destination plugin to serve.
	SourcePlugin      *plugins.SourcePlugin
	DestinationPlugin plugins.DestinationPlugin
	SentryDsn         string
}

const (
	serveShort = `Start plugin server`
	// bufSize used for unit testing grpc server and client
	testBufSize  = 1024 * 1024
	flushTimeout = 5 * time.Second
)

// lis used for unit testing grpc server and client
var testSourceListener *bufconn.Listener
var testDestinationListener *bufconn.Listener
