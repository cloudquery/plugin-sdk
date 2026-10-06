package plugin

import (
	"context"
	"fmt"

	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AssessTables(ctx context.Context, req *pb.AssessTables_Request) (*pb.AssessTables_Response, error) {
	tables := make([]plugin.TablePair, len(req.Tables))
	for i, pair := range req.Tables {
		if len(pair.OldTable) == 0 && len(pair.NewTable) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "table pair %d has neither an old nor a new table", i)
		}
		var err error
		if tables[i].Old, err = tableFromBytes(pair.OldTable); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to decode old table: %v", err)
		}
		if tables[i].New, err = tableFromBytes(pair.NewTable); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to decode new table: %v", err)
		}
	}
	findings, err := s.Plugin.AssessTables(ctx, tables, plugin.AssessOptions{MigrateForce: req.MigrateForce})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assess tables: %v", err)
	}
	resp := &pb.AssessTables_Response{Tables: make([]*pb.AssessTables_TableFinding, len(findings))}
	for i, f := range findings {
		resp.Tables[i] = tableFindingToPB(f)
	}
	return resp, nil
}

func tableFromBytes(b []byte) (*schema.Table, error) {
	if len(b) == 0 {
		return nil, nil
	}
	sc, err := pb.NewSchemaFromBytes(b)
	if err != nil {
		return nil, err
	}
	table, err := schema.NewTableFromArrowSchema(sc)
	if err != nil {
		return nil, fmt.Errorf("failed to create table from schema: %w", err)
	}
	return table, nil
}

func tableFindingToPB(f plugin.TableFinding) *pb.AssessTables_TableFinding {
	columns := make([]*pb.AssessTables_ColumnFinding, len(f.Columns))
	for i, c := range f.Columns {
		columns[i] = &pb.AssessTables_ColumnFinding{
			ColumnName:         c.ColumnName,
			Category:           pb.AssessTables_Category(c.Category),
			OldType:            c.OldType,
			NewType:            c.NewType,
			SafeModeBehavior:   c.SafeModeBehavior,
			ForcedModeBehavior: c.ForcedModeBehavior,
			Evidence:           evidenceToPB(c.Evidence),
		}
	}
	return &pb.AssessTables_TableFinding{
		TableName:                f.TableName,
		Category:                 pb.AssessTables_Category(f.Category),
		SafeModeBehavior:         f.SafeModeBehavior,
		ForcedModeBehavior:       f.ForcedModeBehavior,
		Columns:                  columns,
		Evidence:                 evidenceToPB(f.Evidence),
		IncompleteCoverageReason: f.IncompleteCoverageReason,
	}
}

func evidenceToPB(evidence []plugin.Evidence) []*pb.AssessTables_Evidence {
	res := make([]*pb.AssessTables_Evidence, len(evidence))
	for i, e := range evidence {
		res[i] = &pb.AssessTables_Evidence{SyntheticValue: e.SyntheticValue, Before: e.Before, After: e.After}
	}
	return res
}
