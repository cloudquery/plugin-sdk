package servers

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/internal/pb"
	"github.com/cloudquery/cq-provider-sdk/plugins"
	"github.com/cloudquery/cq-provider-sdk/schema"
	"github.com/cloudquery/cq-provider-sdk/spec"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

type SourceServer struct {
	pb.UnimplementedSourceServer
	Plugin *plugins.SourcePlugin
}

func (s *SourceServer) GetTables(context.Context, *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	b, err := msgpack.Marshal(s.Plugin.Tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tables")
	}
	return &pb.GetTables_Response{
		Tables: b,
	}, nil
}

func (s *SourceServer) GetExampleConfig(context.Context, *pb.GetExampleConfig_Request) (*pb.GetExampleConfig_Response, error) {
	return &pb.GetExampleConfig_Response{
		Name:    s.Plugin.Name,
		Version: s.Plugin.Version,
		Config:  s.Plugin.ExampleConfig}, nil
}

func (s *SourceServer) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec spec.SourceSpec
	if err := yaml.Unmarshal(req.Config, &spec); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	jsonschemaResult, err := s.Plugin.Init(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure source")
	}
	b, err := msgpack.Marshal(jsonschemaResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json schema result")
	}
	return &pb.Configure_Response{
		JsonschemaResult: b,
	}, nil
}

func (s *SourceServer) Fetch(req *pb.Fetch_Request, stream pb.Source_FetchServer) error {
	resources := make(chan *schema.Resource)
	var fetchErr error
	go func() {
		defer close(resources)
		if err := s.Plugin.Fetch(stream.Context(), resources); err != nil {
			fetchErr = errors.Wrap(err, "failed to fetch resources")
		}
	}()

	for resource := range resources {
		b, err := msgpack.Marshal(resource)
		if err != nil {
			return errors.Wrap(err, "failed to marshal resource")
		}
		stream.Send(&pb.Fetch_Response{
			Resource: b,
		})
	}
	if fetchErr != nil {
		return fetchErr
	}

	return nil
}
