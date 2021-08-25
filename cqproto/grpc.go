package cqproto

import (
	"context"

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
		Migrations:     res.Migrations,
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
		Config:        request.Config,
		DisableDelete: request.DisableDelete,
		ExtraFields:   fieldsData,
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
	return &FetchResourcesResponse{
		FinishedResources:           resp.GetFinishedResources(),
		ResourceCount:               resp.GetResourceCount(),
		Error:                       resp.GetError(),
		PartialFetchFailedResources: partialFetchFailedResourcesFromProto(resp.GetPartialFetchFailedResources()),
	}, nil
}

type GRPCServer struct {
	// This is the real implementation
	Impl CQProviderServer
	internal.UnimplementedProviderServer
}

func (g *GRPCServer) GetProviderSchema(ctx context.Context, request *internal.GetProviderSchema_Request) (*internal.GetProviderSchema_Response, error) {
	resp, err := g.Impl.GetProviderSchema(ctx, &GetProviderSchemaRequest{})
	if err != nil {
		return nil, err
	}
	return &internal.GetProviderSchema_Response{
		Name:           resp.Name,
		Version:        resp.Version,
		ResourceTables: tablesToProto(resp.ResourceTables),
		Migrations:     resp.Migrations,
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
		Config:        request.Config,
		DisableDelete: request.DisableDelete,
		ExtraFields:   eFields,
	})
	if err != nil {
		return nil, err
	}
	return &internal.ConfigureProvider_Response{Error: resp.Error}, nil

}

func (g *GRPCServer) FetchResources(request *internal.FetchResources_Request, server internal.Provider_FetchResourcesServer) error {
	return g.Impl.FetchResources(
		server.Context(),
		&FetchResourcesRequest{Resources: request.GetResources(), PartialFetchingEnabled: request.PartialFetchingEnabled},
		&GRPCFetchResourcesServer{server: server},
	)
}

type GRPCFetchResourcesServer struct {
	server internal.Provider_FetchResourcesServer
}

func (g GRPCFetchResourcesServer) Send(response *FetchResourcesResponse) error {
	return g.server.Send(&internal.FetchResources_Response{
		FinishedResources:           response.FinishedResources,
		ResourceCount:               response.ResourceCount,
		Error:                       response.Error,
		PartialFetchFailedResources: partialFetchFailedResourcesToProto(response.PartialFetchFailedResources),
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
		cols[i] = schema.Column{
			Name:        c.GetName(),
			Type:        schema.ValueType(c.GetType()),
			Description: c.GetDescription(),
		}
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

func partialFetchFailedResourcesFromProto(in []*internal.PartialFetchFailedResource) []*PartialFetchFailedResource {
	if len(in) == 0 {
		return nil
	}
	failedResources := make([]*PartialFetchFailedResource, len(in))
	for i, p := range in {
		failedResources[i] = &PartialFetchFailedResource{
			TableName:            p.TableName,
			RootTableName:        p.RootTableName,
			RootPrimaryKeyValues: p.RootPrimaryKeyValues,
			Error:                p.Error,
		}
	}
	return failedResources
}

func partialFetchFailedResourcesToProto(in []*PartialFetchFailedResource) []*internal.PartialFetchFailedResource {
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

// PartialFetchToCQProto converts schema partial fetch failed resources to cq-proto partial fetch resources
func PartialFetchToCQProto(in []schema.PartialFetchFailedResource) []*PartialFetchFailedResource {
	if len(in) == 0 {
		return nil
	}
	failedResources := make([]*PartialFetchFailedResource, len(in))
	for i, p := range in {
		failedResources[i] = &PartialFetchFailedResource{
			TableName:            p.TableName,
			RootTableName:        p.RootTableName,
			RootPrimaryKeyValues: p.RootPrimaryKeyValues,
			Error:                p.Error,
		}
	}
	return failedResources
}
