//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative internal/plugin.proto
package cqproto

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

type CQProvider interface {
	// GetProviderSchema is called when CloudQuery needs to know what the
	// provider's tables and version
	GetProviderSchema(context.Context, *GetProviderSchemaRequest) (*GetProviderSchemaResponse, error)

	// GetProviderConfig is called when CloudQuery wants to generate a configuration example for a provider
	GetProviderConfig(context.Context, *GetProviderConfigRequest) (*GetProviderConfigResponse, error)

	// ConfigureProvider is called to pass the user-specified provider
	// configuration to the provider.
	ConfigureProvider(context.Context, *ConfigureProviderRequest) (*ConfigureProviderResponse, error)

	// FetchResources is called when CloudQuery requests to fetch one or more resources from the provider
	// The provider reports back status updates on the resources fetching progress.
	FetchResources(context.Context, *FetchResourcesRequest) (FetchResourcesStream, error)
}

type CQProviderServer interface {
	// GetProviderSchema is called when CloudQuery needs to know what the
	// provider's tables and version
	GetProviderSchema(context.Context, *GetProviderSchemaRequest) (*GetProviderSchemaResponse, error)

	// GetProviderConfig is called when CloudQuery wants to generate a configuration example for a provider
	GetProviderConfig(context.Context, *GetProviderConfigRequest) (*GetProviderConfigResponse, error)

	// ConfigureProvider is called to pass the user-specified provider
	// configuration to the provider.
	ConfigureProvider(context.Context, *ConfigureProviderRequest) (*ConfigureProviderResponse, error)

	// FetchResources is called when CloudQuery requests to fetch one or more resources from the provider
	// The provider reports back status updates on the resources fetching progress.
	FetchResources(context.Context, *FetchResourcesRequest, FetchResourcesSender) error
}

// GetProviderSchemaRequest represents a CloudQuery RPC request for provider's schemas
type GetProviderSchemaRequest struct{}

type GetProviderSchemaResponse struct {
	// Name is the name of the provider being executed
	Name string
	// Version is the current version provider being executed
	Version string
	// ResourceTables is a map of tables this provider creates
	ResourceTables map[string]*schema.Table
}

// GetProviderConfigRequest represents a CloudQuery RPC request for provider's config
type GetProviderConfigRequest struct{}

type GetProviderConfigResponse struct {
	Config []byte
}

type ConfigureProviderRequest struct {
	// CloudQueryVersion is the version of CloudQuery executing the request.
	CloudQueryVersion string
	// ConnectionDetails holds information regarding connection to the CloudQuery database
	Connection ConnectionDetails
	// Config is the configuration the user supplied for the provider
	Config []byte
	// DisableDelete configures providers to skip deletion of data before resource fetch
	DisableDelete bool
	// Fields to inject to every resource on insert
	ExtraFields map[string]interface{}
}

type ConfigureProviderResponse struct {
	// Error should be set to a string describing the error.
	// The error can be either from malformed configuration or failure to setup
	Error string
}

// FetchResourcesRequest represents a CloudQuery RPC request of one or more resources
type FetchResourcesRequest struct {
	Resources []string
}

// FetchResourcesStream represents a CloudQuery RPC stream of fetch updates from the provider
type FetchResourcesStream interface {
	Recv() (*FetchResourcesResponse, error)
}

// FetchResourcesSender represents a CloudQuery RPC stream of fetch updates from the provider
type FetchResourcesSender interface {
	Send(*FetchResourcesResponse) error
}

// FetchResourcesResponse represents a CloudQuery RPC response of the current fetch progress of the provider
type FetchResourcesResponse struct {
	// map of resources that have finished fetching
	FinishedResources map[string]bool
	// Amount of resources collected so far
	ResourceCount uint64
	// Error value if any, if returned the stream will be canceled
	Error string
}

type ConnectionDetails struct {
	Type string
	DSN  string
}
