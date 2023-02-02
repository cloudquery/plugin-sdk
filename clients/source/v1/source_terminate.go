//go:build !windows

package source

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func (c *Client) terminateProcess() error {
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		c.logger.Error().Err(err).Msg("failed to send interrupt signal to source plugin")
	}
	timer := time.AfterFunc(5*time.Second, func() {
		if err := c.cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill source plugin")
		}
	})
	st, err := c.cmd.Process.Wait()
	timer.Stop()
	if err != nil {
		return err
	}
	if !st.Success() {
		var additionalInfo string
		status := st.Sys().(syscall.WaitStatus)
		if status.Signaled() && st.ExitCode() != -1 {
			additionalInfo += fmt.Sprintf(" (exit code: %d)", st.ExitCode())
		}
		if st.ExitCode() == 137 {
			additionalInfo = " (Out of Memory)"
		}
		return fmt.Errorf("destination plugin process failed with %s%s", st.String(), additionalInfo)
	}

	return nil
}
