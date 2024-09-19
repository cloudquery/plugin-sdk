package opts

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
)

// SchedulerOpts converts plugin.SyncOptions to []scheduler.SyncOption, adding additionalOpts.
func SchedulerOpts(o plugin.SyncOptions, additionalOpts ...scheduler.SyncOption) []scheduler.SyncOption {
	opts := []scheduler.SyncOption{
		scheduler.WithSyncDeterministicCQID(o.DeterministicCQID),
	}
	if o.Shard != nil {
		opts = append(opts, scheduler.WithShard(o.Shard.Num, o.Shard.Total))
	}
	return append(opts, additionalOpts...)
}
