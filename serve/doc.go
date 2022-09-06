package serve

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/docs"
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

			return docs.GenerateSourcePluginDocs(opts.SourcePlugin, args[0])
		},
	}
}
