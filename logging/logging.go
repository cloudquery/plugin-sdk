package logging

import (
	"github.com/hashicorp/go-hclog"

	"os"
)

// New creates a new hclog logger
func New(options *hclog.LoggerOptions) hclog.Logger {
	if options.Level == hclog.NoLevel {

		if options == nil {
			options = &hclog.LoggerOptions{}
		}
		options.Level = hclog.Info
	}
	if options.Output == nil {
		options.Output = os.Stderr
	}
	return hclog.New(options)
}
