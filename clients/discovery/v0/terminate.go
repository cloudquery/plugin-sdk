//go:build !windows

package discovery

import (
	"fmt"
	"os"
	"time"
)

func (c *Client) terminateProcess() error {
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		c.logger.Error().Err(err).Msg("failed to send interrupt signal to plugin")
	}
	timer := time.AfterFunc(5*time.Second, func() {
		if err := c.cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill plugin")
		}
	})
	st, err := c.cmd.Process.Wait()
	timer.Stop()
	if err != nil {
		return err
	}
	if !st.Success() {
		return fmt.Errorf("plugin process exited with status %s", st.String())
	}

	return nil
}
