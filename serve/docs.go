package serve

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/docs"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
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

func (s *PluginServe) newCmdPluginDoc() *cobra.Command {
	format := newEnum([]string{"json", "markdown"}, "markdown")
	cmd := &cobra.Command{
		Use:   "doc <directory>",
		Short: pluginDocShort,
		Long:  pluginDocLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.plugin.Init(cmd.Context(), nil, plugin.NewClientOptions{
				NoConnection: true,
			}); err != nil {
				return err
			}
			tables, err := s.plugin.Tables(cmd.Context(), plugin.TableOptions{
				Tables: []string{"*"},
			})
			if err != nil {
				return err
			}
			g := docs.NewGenerator(s.plugin.Name(), tables)
			if format.Value != "json" {
				return errors.New("only json format is supported. If need to generate markdown, use the `cloudquery tables` command")
			}

			return g.GenerateJSON(args[0], docs.FormatJSON)
		},
	}
	cmd.Flags().Var(format, "format", fmt.Sprintf("output format. one of: %s", strings.Join(format.Allowed, ",")))
	return cmd
}
