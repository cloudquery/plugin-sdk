package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	pluginInfoShort = "Print build information about this plugin"
	pluginInfoLong  = "Print build information about this plugin"
)

func (s *PluginServe) newCmdPluginInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: pluginInfoShort,
		Long:  pluginInfoLong,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println(fmt.Sprintf("Package and version: %s/%s/%s@%s", s.plugin.Team(), s.plugin.Kind(), s.plugin.Name(), s.plugin.Version()))
			return nil
		},
	}
	return cmd
}
