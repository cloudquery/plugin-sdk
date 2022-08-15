package servers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type SourceServer struct {
	pb.UnimplementedSourceServer
	Plugin *plugins.SourcePlugin
}

func (s *SourceServer) GetTables(context.Context, *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	b, err := json.Marshal(s.Plugin.Tables())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tables")
	}
	return &pb.GetTables_Response{
		Tables: b,
	}, nil
}

func (s *SourceServer) GetExampleConfig(context.Context, *pb.GetExampleConfig_Request) (*pb.GetExampleConfig_Response, error) {
	return &pb.GetExampleConfig_Response{
		Name:    s.Plugin.Name(),
		Version: s.Plugin.Version(),
		Config:  s.Plugin.ExampleConfig()}, nil
}

func (s *SourceServer) Fetch(req *pb.Fetch_Request, stream pb.Source_FetchServer) error {
	resources := make(chan *schema.Resource)
	var fetchErr error

	var spec specs.SourceSpec
	if err := yaml.Unmarshal(req.Spec, &spec); err != nil {
		return fmt.Errorf("failed to unmarshal source spec: %w", err)
	}

	go func() {
		defer close(resources)
		if err := s.Plugin.Sync(stream.Context(), spec, resources); err != nil {
			fetchErr = errors.Wrap(err, "failed to fetch resources")
		}
	}()

	for resource := range resources {
		b, err := json.Marshal(schema.WireResource{
			Data:      resource.Data,
			TableName: resource.Table.Name,
		})
		if err != nil {
			return errors.Wrap(err, "failed to marshal resource")
		}
		if err := stream.Send(&pb.Fetch_Response{
			Resource: b,
		}); err != nil {
			return errors.Wrap(err, "failed to send resource")
		}
	}
	if fetchErr != nil {
		return fetchErr
	}

	return nil
}
