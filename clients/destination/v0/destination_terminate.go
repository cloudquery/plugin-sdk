//go:build !windows

package destination

import (
	"fmt"
	"os"
	"time"
)

func (c *Client) terminateProcess() error {
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		c.logger.Error().Err(err).Msg("failed to send interrupt signal to destination plugin")
	}
	timer := time.AfterFunc(5*time.Second, func() {
		c.logger.Info().Msg("sending kill signal to destination plugin")
		if err := c.cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill destination plugin")
		}
	})
	c.logger.Info().Msg("waiting for destination plugin to terminate")
	st, err := c.cmd.Process.Wait()
	timer.Stop()
	if err != nil {
		return err
	}
	if !st.Success() {
		return fmt.Errorf("destination plugin process exited with status %s", st.String())
	}

	return nil
}
