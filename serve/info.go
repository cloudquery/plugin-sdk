package serve

import (
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
			cmd.Println("Package and version:", s.plugin.PackageAndVersion())
			return nil
		},
	}
	return cmd
}
