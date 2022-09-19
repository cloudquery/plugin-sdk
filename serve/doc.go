// Package serve defines APIs to serve (invoke) source and destination plugins
package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	docShort = "Generate markdown documentation for tables"
)

func newCmdDoc(opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "doc <folder>",
		Short: docShort,
		Long:  docShort,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.SourcePlugin == nil {
				return fmt.Errorf("doc generation is only supported for source plugins")
			}
			return opts.SourcePlugin.GenerateSourcePluginDocs(args[0])
		},
	}
}
