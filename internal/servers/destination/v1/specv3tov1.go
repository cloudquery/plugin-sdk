package destination

import (
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-pb-go/specs"
)

func SourceSpecV1ToV3(spec specs.Source) pbPlugin.Spec {
	newSpec := pbPlugin.Spec{
		Name: spec.Name,
		Version: spec.Version,
		Path: spec.Path,
		SyncSpec: &pbPlugin.SyncSpec{
			Tables: spec.Tables,
			SkipTables: spec.SkipTables,
			Destinations: spec.Destinations,
			Concurrency: uint64(spec.Concurrency),
			DetrministicCqId: spec.DeterministicCQID,
		},
	}
	switch spec.Scheduler {
	case specs.SchedulerDFS:
			newSpec.SyncSpec.Scheduler = pbPlugin.SyncSpec_SCHEDULER_DFS
	case specs.SchedulerRoundRobin:
			newSpec.SyncSpec.Scheduler = pbPlugin.SyncSpec_SCHEDULER_ROUND_ROBIN
	default:
		panic("invalid scheduler " + spec.Scheduler.String())
	}
	return newSpec
}

func SpecV1ToV3(spec specs.Destination) pbPlugin.Spec {
	newSpec := pbPlugin.Spec{
		Name: spec.Name,
		Version: spec.Version,
		Path: spec.Path,
		WriteSpec: &pbPlugin.WriteSpec{
			BatchSize: uint64(spec.BatchSize),
			BatchSizeBytes: uint64(spec.BatchSizeBytes),
		},
	}
	switch spec.Registry {
	case specs.RegistryGithub:
			newSpec.Registry = pbPlugin.Spec_REGISTRY_GITHUB
	case specs.RegistryGrpc:
			newSpec.Registry = pbPlugin.Spec_REGISTRY_GRPC
	case specs.RegistryLocal:
			newSpec.Registry = pbPlugin.Spec_REGISTRY_LOCAL
	default:
		panic("invalid registry " + spec.Registry.String())
	}
	switch spec.WriteMode {
	case specs.WriteModeAppend:
			newSpec.WriteSpec.WriteMode = pbPlugin.WRITE_MODE_WRITE_MODE_APPEND
	case specs.WriteModeOverwrite:
			newSpec.WriteSpec.WriteMode = pbPlugin.WRITE_MODE_WRITE_MODE_OVERWRITE
	case specs.WriteModeOverwriteDeleteStale:
			newSpec.WriteSpec.WriteMode = pbPlugin.WRITE_MODE_WRITE_MODE_OVERWRITE_DELETE_STALE
	default:
		panic("invalid write mode " + spec.WriteMode.String())
	}
	switch spec.PKMode {
	case specs.PKModeDefaultKeys:
			newSpec.WriteSpec.PkMode = pbPlugin.WriteSpec_DEFAULT
	case specs.PKModeCQID:
			newSpec.WriteSpec.PkMode = pbPlugin.WriteSpec_CQ_ID_ONLY
	}
	switch spec.MigrateMode {
	case specs.MigrateModeSafe:
			newSpec.WriteSpec.MigrateMode = pbPlugin.WriteSpec_SAFE
	case specs.MigrateModeForced:
			newSpec.WriteSpec.MigrateMode = pbPlugin.WriteSpec_FORCE
	default:
		panic("invalid migrate mode " + spec.MigrateMode.String())
	}
	return newSpec
}