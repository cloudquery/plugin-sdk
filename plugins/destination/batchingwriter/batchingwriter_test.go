package batchingwriter

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type client struct {
	spec specs.Destination
	destination.UnimplementedClient
}

func newClient(_ context.Context, _ zerolog.Logger, spec specs.Destination) (destination.Client, error) {
	return &client{
		spec: spec,
	}, nil
}

func TestBatchSizeInit(t *testing.T) {
	const (
		batchSize      = 100
		batchSizeBytes = 1000
	)

	var (
		batchSizeObserved      int
		batchSizeBytesObserved int
	)
	p := destination.NewPlugin(
		"test",
		"development",
		func(ctx context.Context, logger zerolog.Logger, s specs.Destination) (destination.Client, error) {
			batchSizeObserved = s.BatchSize
			batchSizeBytesObserved = s.BatchSizeBytes
			return newClient(ctx, logger, s)
		},
		destination.WithManagedWriter(New(
			WithDefaultBatchSize(batchSize, batchSizeBytes),
		)),
	)
	require.NoError(t, p.Init(context.TODO(), zerolog.Nop(), specs.Destination{}))

	require.Equal(t, batchSize, batchSizeObserved)
	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
	assertBatchSizes(t, p, int64(batchSize), int64(batchSizeBytes))
}

func TestPluginInitWithSpec(t *testing.T) {
	const (
		batchSize      = 100
		batchSizeBytes = 1000
	)

	var (
		batchSizeObserved      int
		batchSizeBytesObserved int
	)
	p := destination.NewPlugin(
		"test",
		"development",
		func(ctx context.Context, logger zerolog.Logger, s specs.Destination) (destination.Client, error) {
			batchSizeObserved = s.BatchSize
			batchSizeBytesObserved = s.BatchSizeBytes
			return newClient(ctx, logger, s)
		},
		destination.WithManagedWriter(New(
			WithDefaultBatchSize(batchSize*4, batchSizeBytes*4), // set arbitrary defaults
		)),
	)
	require.NoError(t, p.Init(context.TODO(), zerolog.Nop(), specs.Destination{
		BatchSize:      batchSize,
		BatchSizeBytes: batchSizeBytes,
	}))

	require.Equal(t, batchSize, batchSizeObserved)
	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
	assertBatchSizes(t, p, int64(batchSize), int64(batchSizeBytes))
}

func TestPluginInitDefaults(t *testing.T) {
	const (
		// defaults from batchingwriter
		batchSize      = 10000
		batchSizeBytes = 5 * 1024 * 1024
	)

	var (
		batchSizeObserved      int
		batchSizeBytesObserved int
	)
	p := destination.NewPlugin(
		"test",
		"development",
		func(ctx context.Context, logger zerolog.Logger, s specs.Destination) (destination.Client, error) {
			batchSizeObserved = s.BatchSize
			batchSizeBytesObserved = s.BatchSizeBytes
			return newClient(ctx, logger, s)
		},
		destination.WithManagedWriter(New()),
	)
	require.NoError(t, p.Init(context.TODO(), zerolog.Nop(), specs.Destination{}))

	require.Equal(t, batchSize, batchSizeObserved)
	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
	assertBatchSizes(t, p, int64(batchSize), int64(batchSizeBytes))
}

func assertBatchSizes(t *testing.T, p *destination.Plugin, expSize, expBytes int64) {
	bw := p.BatchingWriter()
	s, b := bw.(interface {
		BatchSize() (int64, int64)
	}).BatchSize()
	assert.Equal(t, expSize, s, "batchSize does not match")
	assert.Equal(t, expBytes, b, "batchSizeBytes does not match")
}
