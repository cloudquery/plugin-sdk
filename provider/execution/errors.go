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

func defaultErrorClassifier(_ schema.ClientMeta, resourceName string, err error) diag.Diagnostics {
	if _, ok := err.(diag.Diagnostic); ok {
		return nil
	}
	if _, ok := err.(diag.Diagnostics); ok {
		return nil
	}
	if strings.Contains(err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propagated back to the user via the CLI
		return FromError(err, diag.WithResourceName(resourceName), diag.WithSummary(fdLimitMessage), diag.WithType(diag.THROTTLE), diag.WithSeverity(diag.WARNING))
	}
	return nil
}

func ClassifyError(err error, opts ...diag.BaseErrorOption) diag.Diagnostics {
	if err != nil && strings.Contains(err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propagated back to the user via the CLI
		opts = append(opts, diag.WithSeverity(diag.WARNING), diag.WithType(diag.THROTTLE), diag.WithSummary("%s", fdLimitMessage))
	}
	return FromError(err, opts...)
}

func FromError(err error, opts ...diag.BaseErrorOption) diag.Diagnostics {
	switch ti := err.(type) {
	case diag.Diagnostics:
		ret := make(diag.Diagnostics, len(ti))
		for i := range ti {
			ret[i] = diag.NewBaseError(ti[i], diag.RESOLVING, append([]diag.BaseErrorOption{diag.WithNoOverwrite()}, opts...)...)
		}
		return ret
	default:
		e := diag.NewBaseError(err, diag.RESOLVING, append([]diag.BaseErrorOption{diag.WithNoOverwrite()}, opts...)...)
		return diag.Diagnostics{e}
	}
}

func WithResource(resource *schema.Resource) diag.BaseErrorOption {
	if resource == nil {
		return diag.WithResourceId(nil)
	}
	return diag.WithResourceId(resource.PrimaryKeyValues())
}
