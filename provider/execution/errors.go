package execution

import (
	"strings"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

const (
	// fdLimitMessage defines the message for when a client isn't able to fetch because the open fd limit is hit
	fdLimitMessage = "try increasing number of available file descriptors via `ulimit -n 10240` or by increasing timeout via provider specific parameters"
)

type ErrorClassifier func(meta schema.ClientMeta, resourceName string, err error) diag.Diagnostics

func ClassifyError(err error, opts ...diag.BaseErrorOption) diag.Diagnostics {
	if err != nil && strings.Contains(err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propagated back to the user via the CLI
		opts = append(opts, diag.WithSeverity(diag.WARNING), diag.WithType(diag.THROTTLE), diag.WithSummary(fdLimitMessage))
	}
	return fromError(err, opts...)
}

func WithResource(resource *schema.Resource) diag.BaseErrorOption {
	if resource == nil {
		return diag.WithResourceId(nil)
	}
	return diag.WithResourceId(resource.PrimaryKeyValues())
}

func fromError(err error, opts ...diag.BaseErrorOption) diag.Diagnostics {
	baseOpts := append([]diag.BaseErrorOption{diag.WithNoOverwrite()}, opts...)
	switch ti := err.(type) {
	case diag.Diagnostics:
		ret := make(diag.Diagnostics, 0, len(ti))
		for _, d := range ti {
			ret = append(ret, diag.NewBaseError(d, diag.RESOLVING, baseOpts...))
		}
		return ret
	default:
		e := diag.NewBaseError(err, diag.RESOLVING, baseOpts...)
		return diag.Diagnostics{e}
	}
}
