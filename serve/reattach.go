package serve

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-plugin"
)

// parse information on reattaching to unmanaged providers out of a
// JSON-encoded environment variable.
func ParseReattachProviders(reattachPath string) (map[string]*plugin.ReattachConfig, error) {
	unmanagedProviders := map[string]*plugin.ReattachConfig{}
	if reattachPath == "" {
		return unmanagedProviders, nil
	}

	// Open our reattach config file
	cfg, err := os.Open(reattachPath)
	if err != nil {
		return unmanagedProviders, fmt.Errorf("failed to open provider reattach config: %w", err)
	}
	defer cfg.Close()

	var m map[string]ReattachConfig
	if err := json.NewDecoder(cfg).Decode(&m); err != nil {
		return unmanagedProviders, fmt.Errorf("invalid format for CQ_REATTACH_PROVIDERS: %w", err)
	}

	for p, c := range m {
		var addr net.Addr
		switch c.Addr.Network {
		case "unix":
			addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
			if err != nil {
				return unmanagedProviders, fmt.Errorf("invalid unix socket path %q for %q: %w", c.Addr.String, p, err)
			}
		case "tcp":
			addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
			if err != nil {
				return unmanagedProviders, fmt.Errorf("invalid TCP address %q for %q: %w", c.Addr.String, p, err)
			}
		default:
			return unmanagedProviders, fmt.Errorf("unknown address type %q for %q", c.Addr.Network, p)
		}
		unmanagedProviders[p] = &plugin.ReattachConfig{
			Protocol: plugin.Protocol(c.Protocol),
			Pid:      c.Pid,
			Test:     c.Test,
			Addr:     addr,
		}
	}
	return unmanagedProviders, nil
}

func saveProviderReattach(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write CQ_REATTACH_PROVIDERS=%s: %w", path, err)
	}
	return nil
}
