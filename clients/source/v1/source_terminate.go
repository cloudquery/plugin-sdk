//go:build !windows

package source

import (
	"fmt"
	"os"
	"time"
)

func (c *Client) terminateProcess() error {
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		c.logger.Error().Err(err).Msg("failed to send interrupt signal to source plugin")
	}
	timer := time.AfterFunc(5*time.Second, func() {
		c.logger.Info().Msg("sending kill signal to source plugin")
		if err := c.cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill source plugin")
		}
	})
	c.logger.Info().Msg("waiting for source plugin to terminate")
	st, err := c.cmd.Process.Wait()
	timer.Stop()
	if err != nil {
		return err
	}
	if !st.Success() {
		return fmt.Errorf("source plugin process exited with status %s", st.String())
	}

	return nil
}
