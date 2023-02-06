//go:build windows

package destination

func (c *Client) terminateProcess() error {
	if err := c.cmd.Process.Kill(); err != nil {
		c.logger.Error().Err(err).Msg("failed to kill destination plugin")
	}
	c.logger.Info().Msg("waiting for destination plugin to terminate")
	st, err := c.cmd.Process.Wait()
	if err != nil {
		return err
	}
	if !st.Success() {
		// on windows there is no way to shutdown gracefully via signal. Maybe we can do it via grpc api?
		// though it is a bit strange to expose api to shutdown a server :thinking?:
		c.logger.Info().Msgf("destination plugin process exited with %s", st.String())
	}

	return nil
}
