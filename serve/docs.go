package serve

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	pluginDocShort = "Generate documentation for tables"
	pluginDocLong  = `Generate documentation for tables

If format is markdown, a destination directory will be created (if necessary) containing markdown files.
Example:
doc ./output 

If format is JSON, a destination directory will be created (if necessary) with a single json file called __tables.json.
Example:
doc --format json .
`
)

func (*PluginServe) newCmdPluginDoc() *cobra.Command {
	format := newEnum([]string{"json", "markdown"}, "markdown")
	cmd := &cobra.Command{
		Use:   "doc <directory> (DEPRECATED)",
		Short: pluginDocShort,
		Long:  pluginDocLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("this command is deprecated, please use the `cloudquery tables` command for similar functionality https://www.cloudquery.io/docs/reference/cli/cloudquery_tables")
		},
	}
	cmd.Flags().Var(format, "format", fmt.Sprintf("output format. one of: %s", strings.Join(format.Allowed, ",")))
	return cmd
}
