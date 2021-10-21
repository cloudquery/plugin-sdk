package provider

import (
	"context"
	"embed"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/cloudquery/cq-provider-sdk/provider/schema/diag"

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
	// Migrations embedded and passed by the provider to upgrade between versions
	Migrations embed.FS
	// ErrorClassifier allows the provider to classify errors it returns table execution, and return diagnostics to the user
	ErrorClassifier func(meta schema.ClientMeta, resource string, err error) []diag.Diagnostic
	// Database connection string
	dbURL string
	// meta is the provider's client created when configure is called
	meta schema.ClientMeta
	// Whether provider should all Delete on every table before fetching
	disableDelete bool
	// Add extra fields to all resources, these fields don't show up in documentation and are used for internal CQ testing.
	extraFields map[string]interface{}
	// databaseCreator creates a database based on requested engine
	databaseCreator func(ctx context.Context, logger hclog.Logger, dbURL string) (schema.Database, error)
}

func (p *Provider) GetProviderSchema(_ context.Context, _ *cqproto.GetProviderSchemaRequest) (*cqproto.GetProviderSchemaResponse, error) {
	m, err := readProviderMigrationFiles(p.Logger, p.Migrations)
	if err != nil {
		return nil, err
	}
	return &cqproto.GetProviderSchemaResponse{
		Name:           p.Name,
		Version:        p.Version,
		ResourceTables: p.ResourceMap,
		Migrations:     m,
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
			// list of resources to fetch
			resources = %s
			// enables partial fetching, allowing for any failures to not stop full resource pull
			enable_partial_fetch = true
		}`, p.Name, p.Config().Example(), helpers.FormatSlice(funk.Keys(p.ResourceMap).([]string)))

	return &cqproto.GetProviderConfigResponse{Config: hclwrite.Format([]byte(data))}, nil
}

func (p *Provider) ConfigureProvider(_ context.Context, request *cqproto.ConfigureProviderRequest) (*cqproto.ConfigureProviderResponse, error) {
	if p.meta != nil {
		return &cqproto.ConfigureProviderResponse{Error: fmt.Sprintf("provider %s was already configured", p.Name)}, nil
	}
	if p.Logger == nil {
		return &cqproto.ConfigureProviderResponse{Error: fmt.Sprintf("provider %s logger not defined, make sure to run it with serve", p.Name)}, nil
	}
	// set database creator
	if p.databaseCreator == nil {
		p.databaseCreator = func(ctx context.Context, logger hclog.Logger, dbURL string) (schema.Database, error) {
			return schema.NewPgDatabase(ctx, logger, dbURL)
		}
	}

	p.disableDelete = request.DisableDelete
	p.extraFields = request.ExtraFields
	p.dbURL = request.Connection.DSN
	providerConfig := p.Config()
	if err := defaults.Set(providerConfig); err != nil {
		return &cqproto.ConfigureProviderResponse{}, err
	}
	// if we received an empty config we notify in log and only use defaults.
	if len(request.Config) == 0 {
		p.Logger.Info("Received empty configuration, using only defaults")
	} else if err := hclsimple.Decode("config.json", request.Config, nil, providerConfig); err != nil {
		p.Logger.Error("Failed to load configuration.", "error", err)
		return &cqproto.ConfigureProviderResponse{}, err
	}

	client, err := p.Configure(p.Logger, providerConfig)
	if err != nil {
		return &cqproto.ConfigureProviderResponse{}, err
	}

	tables := make(map[string]string)
	for r, t := range p.ResourceMap {
		if err := getTableDuplicates(r, t, tables); err != nil {
			return &cqproto.ConfigureProviderResponse{}, err
		}
	}

	p.meta = client
	return &cqproto.ConfigureProviderResponse{}, nil
}

func (p *Provider) FetchResources(ctx context.Context, request *cqproto.FetchResourcesRequest, sender cqproto.FetchResourcesSender) error {

	if p.meta == nil {
		return fmt.Errorf("provider client is not configured, call ConfigureProvider first")
	}

	if helpers.HasDuplicates(request.Resources) {
		return fmt.Errorf("provider has duplicate resources requested")
	}

	// if resources ["*"] is requested we will fetch all resources
	resources, err := p.interpolateAllResources(request.Resources)
	if err != nil {
		return err
	}

	conn, err := p.databaseCreator(ctx, p.Logger, p.dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database. %w", err)
	}

	defer conn.Close()

	g, gctx := errgroup.WithContext(ctx)
	finishedResources := make(map[string]bool, len(resources))
	l := sync.Mutex{}
	var totalResourceCount uint64 = 0
	for _, resource := range resources {
		table, ok := p.ResourceMap[resource]
		if !ok {
			return fmt.Errorf("plugin %s does not provide resource %s", p.Name, resource)
		}
		execData := schema.NewExecutionData(conn, p.Logger, table, p.disableDelete, p.extraFields, request.PartialFetchingEnabled)
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
			defer l.Unlock()
			if err != nil {
				status := cqproto.ResourceFetchFailed
				if err == context.Canceled {
					status = cqproto.ResourceFetchCanceled
				}
				return sender.Send(&cqproto.FetchResourcesResponse{
					ResourceName:                r,
					FinishedResources:           finishedResources,
					ResourceCount:               resourceCount,
					Error:                       err.Error(),
					PartialFetchFailedResources: cqproto.PartialFetchToCQProto(execData.PartialFetchFailureResult),
					Summary: cqproto.ResourceFetchSummary{
						Status:        status,
						ResourceCount: resourceCount,
						Diagnostics:   p.collectExecutionDiagnostics(p.meta, execData),
					},
				})
			}
			status := cqproto.ResourceFetchComplete
			if len(execData.PartialFetchFailureResult) > 0 {
				status = cqproto.ResourceFetchPartial
			}
			err = sender.Send(&cqproto.FetchResourcesResponse{
				ResourceName:                r,
				FinishedResources:           finishedResources,
				ResourceCount:               resourceCount,
				Error:                       "",
				PartialFetchFailedResources: cqproto.PartialFetchToCQProto(execData.PartialFetchFailureResult),
				Summary: cqproto.ResourceFetchSummary{
					Status:        status,
					ResourceCount: resourceCount,
					Diagnostics:   p.collectExecutionDiagnostics(p.meta, execData),
				},
			})
			if err != nil {
				return err
			}
			p.Logger.Debug("finished fetching table...", "provider", p.Name, "table", table.Name)
			return nil
		})
	}
	return g.Wait()
}

func (p *Provider) collectExecutionDiagnostics(client schema.ClientMeta, exec schema.ExecutionData) diag.Diagnostics {
	classifier := DefaultErrorClassifier
	if p.ErrorClassifier != nil {
		classifier = p.ErrorClassifier
	}
	p.Logger.Debug("collecting diagnostics for resource execution", "resource", exec.ResourceName)
	diagnostics := make(diag.Diagnostics, 0)
	for _, e := range exec.PartialFetchFailureResult {
		if d, ok := e.Err.(diag.Diagnostic); ok {
			diagnostics = append(diagnostics, d)
			continue
		}
		dd := classifier(client, exec.ResourceName, e.Err)
		if len(dd) > 0 {
			diagnostics = append(diagnostics, dd...)
			continue
		}
		// if error wasn't classified by provider mark it as error
		diagnostics = append(diagnostics, diag.FromError(e.Err, diag.ERROR, diag.RESOLVING, exec.ResourceName, e.Error(), ""))
	}
	return diagnostics
}

func (p *Provider) interpolateAllResources(requestedResources []string) ([]string, error) {
	if len(requestedResources) != 1 {
		if funk.ContainsString(requestedResources, "*") {
			return nil, fmt.Errorf("invalid \"*\" resource, with explicit resources")
		}
		return requestedResources, nil
	}
	if requestedResources[0] != "*" {
		return requestedResources, nil
	}
	allResources := make([]string, 0, len(p.ResourceMap))
	for k := range p.ResourceMap {
		allResources = append(allResources, k)
	}
	return allResources, nil
}

func getTableDuplicates(resource string, table *schema.Table, tableNames map[string]string) error {
	for _, r := range table.Relations {
		if err := getTableDuplicates(resource, r, tableNames); err != nil {
			return err
		}
	}
	if existing, ok := tableNames[table.Name]; ok {
		return fmt.Errorf("table name %s used more than once, duplicates are in %s and %s", table.Name, existing, resource)
	}
	tableNames[table.Name] = resource
	return nil
}
