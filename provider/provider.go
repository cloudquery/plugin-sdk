package provider

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/thoas/go-funk"

	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/cloudquery/cq-provider-sdk/helpers"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"golang.org/x/sync/errgroup"
)

// Config Every provider implements a resources field we only want to extract that in fetch execution
type Config interface {
	// Example returns a configuration example (with comments) so user clients can generate an example config
	Example() string
}

// Provider is the base structure required to pass and serve an sdk provider.Provider
type Provider struct {
	// Name of plugin i.e aws,gcp, azure etc'
	Name string
	// Version of the provider
	Version string
	// Configure the provider and return context
	Configure func(hclog.Logger, interface{}) (schema.ClientMeta, error)
	// ResourceMap is all resources supported by this plugin
	ResourceMap map[string]*schema.Table
	// Configuration decoded from configure request
	Config func() Config
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	Logger hclog.Logger
	// Database connection
	db schema.Database
	// meta is the provider's client created when configure is called
	meta schema.ClientMeta
}

func (p *Provider) GetProviderSchema(_ context.Context, _ *cqproto.GetProviderSchemaRequest) (*cqproto.GetProviderSchemaResponse, error) {
	return &cqproto.GetProviderSchemaResponse{
		Name:           p.Name,
		Version:        p.Version,
		ResourceTables: p.ResourceMap,
	}, nil
}

func (p *Provider) GetProviderConfig(_ context.Context, _ *cqproto.GetProviderConfigRequest) (*cqproto.GetProviderConfigResponse, error) {
	providerConfig := p.Config()
	if err := defaults.Set(providerConfig); err != nil {
		return &cqproto.GetProviderConfigResponse{}, err
	}
	data := fmt.Sprintf(`
		provider "%s" {
			%s
			resources = %s
		}`, p.Name, p.Config().Example(), helpers.FormatSlice(funk.Keys(p.ResourceMap).([]string)))
	return &cqproto.GetProviderConfigResponse{Config: hclwrite.Format([]byte(data))}, nil
}

func (p *Provider) ConfigureProvider(_ context.Context, request *cqproto.ConfigureProviderRequest) (*cqproto.ConfigureProviderResponse, error) {
	conn, err := schema.NewPgDatabase(request.Connection.DSN)
	if err != nil {
		return nil, err
	}
	p.db = conn
	// Create tables
	m := NewMigrator(p.db, p.Logger)
	for _, t := range p.ResourceMap {
		// validate table
		if err := schema.ValidateTable(t); err != nil {
			p.Logger.Error("table validation failed", "table", t.Name, "error", err)
			return &cqproto.ConfigureProviderResponse{}, err
		}

		if err := m.CreateTable(context.Background(), t, nil); err != nil {
			p.Logger.Error("failed to create table", "table", t.Name, "error", err)
			return &cqproto.ConfigureProviderResponse{}, err
		}
	}

	providerConfig := p.Config()
	if err := defaults.Set(providerConfig); err != nil {
		return &cqproto.ConfigureProviderResponse{}, err
	}
	if err := hclsimple.Decode("config.json", request.Config, nil, providerConfig); err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	client, err := p.Configure(p.Logger, providerConfig)
	if err != nil {
		return &cqproto.ConfigureProviderResponse{}, err
	}
	p.meta = client
	return &cqproto.ConfigureProviderResponse{}, nil
}

func (p *Provider) FetchResources(ctx context.Context, request *cqproto.FetchResourcesRequest, sender cqproto.FetchResourcesSender) error {

	if p.meta == nil {
		return fmt.Errorf("provider client is not configured, call ConfigureProvider first")
	}

	g, gctx := errgroup.WithContext(ctx)
	finishedResources := make(map[string]bool, len(request.Resources))
	l := sync.Mutex{}
	var totalResourceCount uint64 = 0
	for _, resource := range request.Resources {
		table, ok := p.ResourceMap[resource]
		if !ok {
			return fmt.Errorf("plugin %s does not provide resource %s", p.Name, resource)
		}
		execData := schema.NewExecutionData(p.db, p.Logger, table)
		p.Logger.Debug("fetching table...", "provider", p.Name, "table", table.Name)
		// Save resource aside
		r := resource
		l.Lock()
		finishedResources[r] = false
		l.Unlock()
		g.Go(func() error {
			resourceCount, err := execData.ResolveTable(gctx, p.meta, nil)
			l.Lock()
			finishedResources[r] = true
			atomic.AddUint64(&totalResourceCount, resourceCount)
			l.Unlock()
			if err != nil {
				return sender.Send(&cqproto.FetchResourcesResponse{
					FinishedResources: finishedResources,
					ResourceCount:     resourceCount,
					Error:             err.Error(),
				})
			}
			err = sender.Send(&cqproto.FetchResourcesResponse{
				FinishedResources: finishedResources,
				ResourceCount:     resourceCount,
				Error:             "",
			})
			if err != nil {
				return err
			}
			p.Logger.Debug("fetching table...", "provider", p.Name, "table", table.Name)
			return nil
		})
	}
	return g.Wait()
}
