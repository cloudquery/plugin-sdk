package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

// Every provider implements a resources field we only want to extract that in fetch execution
type Config struct {
	// global timeout in seconds
	Timeout   int `yaml:"timeout" default:"1200"`
	Resources []struct {
		Name  string
		Other map[string]interface{} `yaml:",inline"`
	}
}

// Provider is the base structure required to pass and serve an sdk provider.Provider
type Provider struct {
	// Name of plugin i.e aws,gcp, azure etc'
	Name string
	// Configure the provider and return context
	Configure func(hclog.Logger, interface{}) (schema.ClientMeta, error)
	// ResourceMap is all resources supported by this plugin
	ResourceMap map[string]*schema.Table
	// Configuration unmarshalled from Fetch request
	Config func() interface{}
	// Logger to call
	Logger hclog.Logger
	// DefaultConfigGenerator generates the default configuration for a client to execute this provider
	DefaultConfigGenerator func() (string, error)
	// Database connection
	db schema.Database
}

func (p *Provider) GenConfig() (string, error) {
	return p.DefaultConfigGenerator()
}

func (p *Provider) Init(_ string, dsn string, _ bool) error {
	if p.Logger == nil {
		p.Logger = logging.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			JSONFormat: true,
		})
	}
	conn, err := schema.NewPgDatabase(dsn)
	if err != nil {
		return err
	}
	p.db = conn
	// Create tables
	m := NewMigrator(p.db, p.Logger)
	for _, t := range p.ResourceMap {
		err := m.CreateTable(context.Background(), t, nil)
		if err != nil {
			p.Logger.Error("failed to create table", "table", t.Name, "error", err)
			return err
		}
	}
	return nil
}

func (p *Provider) Fetch(data []byte) error {
	providerConfig := p.Config()
	if err := defaults.Set(providerConfig); err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, providerConfig); err != nil {
		return err
	}

	providerClient, err := p.Configure(p.Logger, providerConfig)
	if err != nil {
		return err
	}

	var providerCfg Config
	if err := defaults.Set(&providerCfg); err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &providerCfg); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(providerCfg.Timeout)*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	for _, r := range providerCfg.Resources {
		table, ok := p.ResourceMap[r.Name]
		if !ok {
			return fmt.Errorf("plugin %s does not provide resource %s", p.Name, r.Name)
		}
		execData := schema.NewExecutionData(p.db, p.Logger, table)
		g.Go(func() error {
			return execData.ResolveTable(ctx, providerClient, nil)
		})
	}
	return g.Wait()
}
