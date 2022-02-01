package cqproto

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/cloudquery/cq-provider-sdk/cqproto/internal"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-plugin"
)

type GRPCClient struct {
	broker *plugin.GRPCBroker
	client internal.ProviderClient
}

func (g GRPCClient) GetProviderSchema(ctx context.Context, _ *GetProviderSchemaRequest) (*GetProviderSchemaResponse, error) {
	res, err := g.client.GetProviderSchema(ctx, &internal.GetProviderSchema_Request{})
	if err != nil {
		return nil, err
	}
	resp := &GetProviderSchemaResponse{
		Name:           res.GetName(),
		Version:        res.GetVersion(),
		ResourceTables: tablesFromProto(res.GetResourceTables()),
		Migrations:     migrationsFromProto(res.GetMigrations()),
	}

	return resp, nil
}

func (g GRPCClient) GetProviderConfig(ctx context.Context, _ *GetProviderConfigRequest) (*GetProviderConfigResponse, error) {
	res, err := g.client.GetProviderConfig(ctx, &internal.GetProviderConfig_Request{})
	if err != nil {
		return nil, err
	}
	return &GetProviderConfigResponse{
		Config: res.GetConfig(),
	}, nil
}

func (g GRPCClient) ConfigureProvider(ctx context.Context, request *ConfigureProviderRequest) (*ConfigureProviderResponse, error) {
	fieldsData, err := msgpack.Marshal(request.ExtraFields)
	if err != nil {
		return nil, err
	}
	res, err := g.client.ConfigureProvider(ctx, &internal.ConfigureProvider_Request{
		CloudqueryVersion: request.CloudQueryVersion,
		Connection: &internal.ConnectionDetails{
			Type: internal.ConnectionType_POSTGRES,
			Dsn:  request.Connection.DSN,
		},
		Config:      request.Config,
		ExtraFields: fieldsData,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureProviderResponse{res.GetError()}, nil
}

func (g GRPCClient) FetchResources(ctx context.Context, request *FetchResourcesRequest) (FetchResourcesStream, error) {
	res, err := g.client.FetchResources(ctx, &internal.FetchResources_Request{
		Resources:              request.Resources,
		PartialFetchingEnabled: request.PartialFetchingEnabled,
		ParallelFetchingLimit:  request.ParallelFetchingLimit,
	})
	if err != nil {
		return nil, err
	}
	return &GRPCFetchResponseStream{res}, nil
}

type GRPCFetchResponseStream struct {
	stream internal.Provider_FetchResourcesClient
}

func (g GRPCFetchResponseStream) Recv() (*FetchResourcesResponse, error) {
	resp, err := g.stream.Recv()
	if err != nil {
		return nil, err
	}
	fr := &FetchResourcesResponse{
		ResourceName:                resp.GetResource(),
		FinishedResources:           resp.GetFinishedResources(),
		ResourceCount:               resp.GetResourceCount(),
		Error:                       resp.GetError(),
		PartialFetchFailedResources: partialFetchFailedResourcesFromProto(resp.GetPartialFetchFailedResources()),
	}
	if resp.GetSummary() != nil {
		fr.Summary = ResourceFetchSummary{
			Status:        ResourceFetchStatus(resp.Summary.Status),
			ResourceCount: resp.GetSummary().GetResourceCount(),
			Diagnostics:   diagnosticsFromProto(resp.GetResource(), resp.GetSummary().Diagnostics),
		}
	}
	return fr, nil
}

type GRPCServer struct {
	// This is the real implementation
	Impl CQProviderServer
	internal.UnimplementedProviderServer
}

func (g *GRPCServer) GetProviderSchema(ctx context.Context, _ *internal.GetProviderSchema_Request) (*internal.GetProviderSchema_Response, error) {
	resp, err := g.Impl.GetProviderSchema(ctx, &GetProviderSchemaRequest{})
	if err != nil {
		return nil, err
	}
	return &internal.GetProviderSchema_Response{
		Name:           resp.Name,
		Version:        resp.Version,
		ResourceTables: tablesToProto(resp.ResourceTables),
		Migrations:     migrationsToProto(resp.Migrations),
	}, nil

}

func (g *GRPCServer) GetProviderConfig(ctx context.Context, _ *internal.GetProviderConfig_Request) (*internal.GetProviderConfig_Response, error) {
	resp, err := g.Impl.GetProviderConfig(ctx, &GetProviderConfigRequest{})
	if err != nil {
		return nil, err
	}
	return &internal.GetProviderConfig_Response{Config: resp.Config}, nil
}

func (g *GRPCServer) ConfigureProvider(ctx context.Context, request *internal.ConfigureProvider_Request) (*internal.ConfigureProvider_Response, error) {

	var eFields = make(map[string]interface{})
	if request.GetExtraFields() != nil {
		if err := msgpack.Unmarshal(request.GetExtraFields(), &eFields); err != nil {
			return nil, err
		}
	}
	resp, err := g.Impl.ConfigureProvider(ctx, &ConfigureProviderRequest{
		CloudQueryVersion: request.GetCloudqueryVersion(),
		Connection: ConnectionDetails{
			Type: string(request.Connection.GetType()),
			DSN:  request.Connection.GetDsn(),
		},
		Config:      request.Config,
		ExtraFields: eFields,
	})
	if err != nil {
		return nil, err
	}
	return &internal.ConfigureProvider_Response{Error: resp.Error}, nil

}

func (g *GRPCServer) FetchResources(request *internal.FetchResources_Request, server internal.Provider_FetchResourcesServer) error {
	return g.Impl.FetchResources(
		server.Context(),
		&FetchResourcesRequest{Resources: request.GetResources(), PartialFetchingEnabled: request.PartialFetchingEnabled, ParallelFetchingLimit: request.ParallelFetchingLimit},
		&GRPCFetchResourcesServer{server: server},
	)
}

type GRPCFetchResourcesServer struct {
	server internal.Provider_FetchResourcesServer
}

func (g GRPCFetchResourcesServer) Send(response *FetchResourcesResponse) error {
	return g.server.Send(&internal.FetchResources_Response{
		Resource:                    response.ResourceName,
		FinishedResources:           response.FinishedResources,
		ResourceCount:               response.ResourceCount,
		Error:                       response.Error,
		PartialFetchFailedResources: partialFetchFailedResourcesToProto(response.PartialFetchFailedResources),
		Summary: &internal.ResourceFetchSummary{
			Status:        internal.ResourceFetchSummary_Status(response.Summary.Status),
			ResourceCount: response.Summary.ResourceCount,
			Diagnostics:   diagnosticsToProto(response.Summary.Diagnostics),
		},
	})
}

func tablesFromProto(in map[string]*internal.Table) map[string]*schema.Table {
	if in == nil {
		return nil
	}
	out := make(map[string]*schema.Table, len(in))
	for k, v := range in {
		out[k] = tableFromProto(v)
	}
	return out
}

func tableFromProto(v *internal.Table) *schema.Table {
	cols := make([]schema.Column, len(v.GetColumns()))
	for i, c := range v.GetColumns() {
		cols[i] = schema.SetColumnMeta(schema.Column{
			Name:        c.GetName(),
			Type:        schema.ValueType(c.GetType()),
			Description: c.GetDescription(),
		}, metaFromProto(c.GetMeta()))
	}
	rels := make([]*schema.Table, len(v.GetRelations()))
	for i, r := range v.GetRelations() {
		rels[i] = tableFromProto(r)
	}

	var opts schema.TableCreationOptions
	if o := v.GetOptions(); o != nil {
		opts.PrimaryKeys = o.GetPrimaryKeys()
	}

	return &schema.Table{
		Name:        v.GetName(),
		Description: v.GetDescription(),
		Columns:     cols,
		Relations:   rels,
		Options:     opts,
	}
}

func metaFromProto(m *internal.ColumnMeta) *schema.ColumnMeta {
	if m == nil {
		return nil
	}
	var r *schema.ResolverMeta
	if m.GetResolver() != nil {
		r = &schema.ResolverMeta{
			Name:    m.Resolver.Name,
			Builtin: m.Resolver.Builtin,
		}
	}
	return &schema.ColumnMeta{
		Resolver:     r,
		IgnoreExists: m.GetIgnoreExists(),
	}
}

func tablesToProto(in map[string]*schema.Table) map[string]*internal.Table {
	if in == nil {
		return nil
	}
	out := make(map[string]*internal.Table, len(in))
	for k, v := range in {
		out[k] = tableToProto(v)
	}
	return out
}

func tableToProto(in *schema.Table) *internal.Table {
	cols := make([]*internal.Column, len(in.Columns))
	for i, c := range in.Columns {
		cols[i] = &internal.Column{
			Name:        c.Name,
			Type:        internal.ColumnType(c.Type),
			Description: c.Description,
			Meta:        columnMetaToProto(c.Meta()),
		}
	}
	rels := make([]*internal.Table, len(in.Relations))
	for i, r := range in.Relations {
		rels[i] = tableToProto(r)
	}
	return &internal.Table{
		Name:        in.Name,
		Description: in.Description,
		Columns:     cols,
		Relations:   rels,
		Options: &internal.TableCreationOptions{
			PrimaryKeys: in.Options.PrimaryKeys,
		},
	}
}

func columnMetaToProto(m *schema.ColumnMeta) *internal.ColumnMeta {
	if m == nil {
		return nil
	}
	var r *internal.ResolverMeta
	if m.Resolver != nil {
		r = &internal.ResolverMeta{Name: m.Resolver.Name, Builtin: m.Resolver.Builtin}
	}
	return &internal.ColumnMeta{
		Resolver:     r,
		IgnoreExists: m.IgnoreExists,
	}
}

func partialFetchFailedResourcesFromProto(in []*internal.PartialFetchFailedResource) []*FailedResourceFetch {
	if len(in) == 0 {
		return nil
	}
	failedResources := make([]*FailedResourceFetch, len(in))
	for i, p := range in {
		failedResources[i] = &FailedResourceFetch{
			TableName:            p.TableName,
			RootTableName:        p.RootTableName,
			RootPrimaryKeyValues: p.RootPrimaryKeyValues,
			Error:                p.Error,
		}
	}
	return failedResources
}

func partialFetchFailedResourcesToProto(in []*FailedResourceFetch) []*internal.PartialFetchFailedResource {
	if len(in) == 0 {
		return nil
	}
	failedResources := make([]*internal.PartialFetchFailedResource, len(in))
	for i, p := range in {
		failedResources[i] = &internal.PartialFetchFailedResource{
			TableName:            p.TableName,
			RootTableName:        p.RootTableName,
			RootPrimaryKeyValues: p.RootPrimaryKeyValues,
			Error:                p.Error,
		}
	}
	return failedResources
}

func diagnosticsToProto(in diag.Diagnostics) []*internal.Diagnostic {
	if len(in) == 0 {
		return nil
	}
	diagnostics := make([]*internal.Diagnostic, len(in))
	for i, p := range in {
		diagnostics[i] = &internal.Diagnostic{
			Type:     internal.Diagnostic_Type(p.Type()),
			Severity: internal.Diagnostic_Severity(p.Severity()),
			Summary:  p.Description().Summary,
			Detail:   p.Description().Detail,
			Resource: p.Description().Resource,
		}
	}
	return diagnostics
}

func diagnosticsFromProto(resourceName string, in []*internal.Diagnostic) diag.Diagnostics {
	if len(in) == 0 {
		return nil
	}
	diagnostics := make(diag.Diagnostics, len(in))
	for i, p := range in {
		diagnostics[i] = &ProviderDiagnostic{
			ResourceName:       resourceName,
			DiagnosticType:     diag.DiagnosticType(p.GetType()),
			DiagnosticSeverity: diag.Severity(p.GetSeverity()),
			Summary:            p.GetSummary(),
			Details:            p.GetDetail(),
		}
	}
	return diagnostics
}

func migrationsFromProto(in map[string]*internal.DialectMigration) map[string]map[string][]byte {
	ret := make(map[string]map[string][]byte, len(in))
	for k := range in {
		ret[k] = in[k].Migrations
	}
	return ret
}

func migrationsToProto(in map[string]map[string][]byte) map[string]*internal.DialectMigration {
	ret := make(map[string]*internal.DialectMigration, len(in))
	for k := range in {
		ret[k] = &internal.DialectMigration{
			Migrations: in[k],
		}
	}
	return ret
}
