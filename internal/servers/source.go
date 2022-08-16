package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/pkg/errors"
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

func (s *SourceServer) ExampleConfig(context.Context, *pb.GetExampleConfig_Request) (*pb.GetExampleConfig_Response, error) {
	exampleConfig, err := s.Plugin.ExampleConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get example config: %w", err)
	}
	return &pb.GetExampleConfig_Response{
		Name:    s.Plugin.Name(),
		Version: s.Plugin.Version(),
		Config:  exampleConfig}, nil
}

func (s *SourceServer) Sync(req *pb.Fetch_Request, stream pb.Source_FetchServer) error {
	resources := make(chan *schema.Resource)
	var fetchErr error

	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode source spec: %w", err)
	}

	go func() {
		defer close(resources)
		if err := s.Plugin.Sync(stream.Context(), spec, resources); err != nil {
			fetchErr = errors.Wrap(err, "failed to fetch resources")
		}
	}()

	for resource := range resources {
		b, err := json.Marshal(resource)
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
