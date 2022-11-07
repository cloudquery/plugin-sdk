package plugins

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type DestinationPluginTestSuite struct {
	skipTestOverwrite            bool
	skipTestOverWriteDeleteStale bool
	skipTestAppend               bool
}

type DestinationPluginTestSuiteOption func(suite *DestinationPluginTestSuite)

func DestinationPluginTestWithSuiteSkipTestOverwrite() DestinationPluginTestSuiteOption {
	return func(suite *DestinationPluginTestSuite) {
		suite.skipTestOverwrite = true
	}
}

func DestinationPluginTestWithSuiteSkipTestOverWriteDeleteStale() DestinationPluginTestSuiteOption {
	return func(suite *DestinationPluginTestSuite) {
		suite.skipTestOverWriteDeleteStale = true
	}
}

func DestinationPluginTestWithSuiteSkipTestAppend() DestinationPluginTestSuiteOption {
	return func(suite *DestinationPluginTestSuite) {
		suite.skipTestAppend = true
	}
}

// TestTable returns a table with columns of all type. useful for destination testing purposes
func testTable(name string) *schema.Table {
	return &schema.Table{
		Name:        name,
		Description: "Test table",
		Columns: schema.ColumnList{
			schema.CqIDColumn,
			schema.CqParentIDColumn,
			{
				Name: "bool",
				Type: schema.TypeBool,
			},
			{
				Name: "int",
				Type: schema.TypeInt,
			},
			{
				Name: "float",
				Type: schema.TypeFloat,
			},
			{
				Name:            "uuid",
				Type:            schema.TypeUUID,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "text",
				Type: schema.TypeString,
			},
			{
				Name: "bytea",
				Type: schema.TypeByteArray,
			},
			{
				Name: "text_array",
				Type: schema.TypeStringArray,
			},
			{
				Name: "int_array",
				Type: schema.TypeIntArray,
			},
			{
				Name: "timestamp",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "json",
				Type: schema.TypeJSON,
			},
			{
				Name: "uuid_array",
				Type: schema.TypeUUIDArray,
			},
			{
				Name: "inet",
				Type: schema.TypeInet,
			},
			{
				Name: "inet_array",
				Type: schema.TypeInetArray,
			},
			{
				Name: "cidr",
				Type: schema.TypeCIDR,
			},
			{
				Name: "cidr_array",
				Type: schema.TypeCIDRArray,
			},
			{
				Name: "macaddr",
				Type: schema.TypeMacAddr,
			},
			{
				Name: "macaddr_array",
				Type: schema.TypeMacAddrArray,
			},
		},
	}
}

func TestData() schema.CQTypes {
	cqID := &schema.UUID{}
	if err := cqID.Set(uuid.New()); err != nil {
		panic(err)
	}
	cqParentID := &schema.UUID{}
	if err := cqParentID.Set("00000000-0000-0000-0000-000000000003"); err != nil {
		panic(err)
	}
	boolColumn := &schema.Bool{
		Bool:   true,
		Status: schema.Present,
	}
	intColumn := &schema.Int8{
		Int:    1,
		Status: schema.Present,
	}
	floatColumn := &schema.Float8{
		Float:  1.1,
		Status: schema.Present,
	}
	uuidColumn := &schema.UUID{}
	if err := uuidColumn.Set("00000000-0000-0000-0000-000000000001"); err != nil {
		panic(err)
	}
	textColumn := &schema.Text{
		Str:    "test",
		Status: schema.Present,
	}
	byteaColumn := &schema.Bytea{
		Bytes:  []byte{1, 2, 3},
		Status: schema.Present,
	}
	textArrayColumn := &schema.TextArray{}
	if err := textArrayColumn.Set([]string{"test1", "test2"}); err != nil {
		panic(err)
	}
	intArrayColumn := &schema.Int8Array{}
	if err := intArrayColumn.Set([]int8{1, 2}); err != nil {
		panic(err)
	}
	timestampColumn := &schema.Timestamptz{
		Time:   time.Now(),
		Status: schema.Present,
	}
	jsonColumn := &schema.JSON{
		Bytes:  []byte(`{"test": "test"}`),
		Status: schema.Present,
	}
	uuidArrayColumn := &schema.UUIDArray{}
	if err := uuidArrayColumn.Set([]string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}); err != nil {
		panic(err)
	}
	inetColumn := &schema.Inet{}
	if err := inetColumn.Set("192.0.2.1/24"); err != nil {
		panic(err)
	}
	inetArrayColumn := &schema.InetArray{}
	if err := inetArrayColumn.Set([]string{"192.0.2.1/24", "192.0.2.1/24"}); err != nil {
		panic(err)
	}
	cidrColumn := &schema.CIDR{}
	if err := cidrColumn.Set("192.0.2.1"); err != nil {
		panic(err)
	}
	cidrArrayColumn := &schema.CIDRArray{}
	if err := cidrArrayColumn.Set([]string{"192.0.2.1", "192.0.2.1"}); err != nil {
		panic(err)
	}
	macaddrColumn := &schema.Macaddr{}
	if err := macaddrColumn.Set("aa:bb:cc:dd:ee:ff"); err != nil {
		panic(err)
	}
	macaddrArrayColumn := &schema.MacaddrArray{}
	if err := macaddrArrayColumn.Set([]string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"}); err != nil {
		panic(err)
	}

	data := schema.CQTypes{
		cqID,
		cqParentID,
		boolColumn,
		intColumn,
		floatColumn,
		uuidColumn,
		textColumn,
		byteaColumn,
		textArrayColumn,
		intArrayColumn,
		timestampColumn,
		jsonColumn,
		uuidArrayColumn,
		inetColumn,
		inetArrayColumn,
		cidrColumn,
		cidrArrayColumn,
		macaddrColumn,
		macaddrArrayColumn,
	}

	return data
}

func getTestLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func (s *DestinationPluginTestSuite) destinationPluginTestWriteOverwrite(ctx context.Context, p *DestinationPlugin, logger zerolog.Logger, spec specs.Destination) error {
	if s.skipTestOverwrite {
		return nil
	}
	spec.WriteMode = specs.WriteModeOverwrite
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_overwrite"
	table := testTable(tableName)
	syncTime := time.Now().UTC()
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "cq_test_write_overwrite_source"
	resource := schema.DestinationResource{
		TableName: table.Name,
		Data:      TestData(),
	}
	resource2 := schema.DestinationResource{
		TableName: table.Name,
		Data:      TestData(),
	}
	_ = resource2.Data[5].Set("00000000-0000-0000-0000-000000000007")
	resources := []schema.DestinationResource{
		resource,
		resource2,
	}

	if !s.skipTestOverWriteDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, syncTime); err != nil {
			return fmt.Errorf("failed to delete stale data: %w", err)
		}
	}

	if err := p.writeAll(ctx, tables, sourceName, syncTime, resources); err != nil {
		return fmt.Errorf("failed to write one: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resource, got %d", len(resourcesRead))
	}

	if resource.Data.Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	if resource2.Data.Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	secondSyncTime := time.Now().UTC()
	// write second time
	if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resource, got %d", len(resourcesRead))
	}

	if resource.Data.Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	if resource2.Data.Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	if !s.skipTestOverWriteDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, secondSyncTime); err != nil {
			return fmt.Errorf("failed to delete stale data second time: %w", err)
		}
	}

	resourcesRead, err = p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}
	if len(resourcesRead) != 1 {
		return fmt.Errorf("expected 1 resource, got %d", len(resourcesRead))
	}

	if resource2.Data.Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	return nil
}

func (*DestinationPluginTestSuite) destinationPluginTestWriteAppend(ctx context.Context, p *DestinationPlugin, logger zerolog.Logger, spec specs.Destination) error {
	spec.WriteMode = specs.WriteModeAppend
	if err := p.Init(ctx, logger, spec); err != nil {
		return fmt.Errorf("failed to init plugin: %w", err)
	}
	tableName := "cq_test_write_append"
	table := testTable(tableName)
	syncTime := time.Now().UTC()
	tables := []*schema.Table{
		table,
	}
	if err := p.Migrate(ctx, tables); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	sourceName := "cq_test_write_overwrite_append"
	resource := schema.DestinationResource{
		TableName: table.Name,
		Data:      TestData(),
	}

	if err := p.writeOne(ctx, tables, sourceName, syncTime, resource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	resource = schema.DestinationResource{
		TableName: table.Name,
		Data:      TestData(),
	}
	secondSyncTime := time.Now().UTC()
	// write second time
	if err := p.writeOne(ctx, tables, sourceName, secondSyncTime, resource); err != nil {
		return fmt.Errorf("failed to write one second time: %w", err)
	}

	if err := p.DeleteStale(ctx, tables, sourceName, syncTime); err != nil {
		return fmt.Errorf("failed to delete stale data: %w", err)
	}

	resourcesRead, err := p.readAll(ctx, tables[0], sourceName)
	if err != nil {
		return fmt.Errorf("failed to read all second time: %w", err)
	}

	if len(resourcesRead) != 2 {
		return fmt.Errorf("expected 2 resource, got %d", len(resourcesRead))
	}

	if resource.Data.Equal(resourcesRead[0]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	if resource.Data.Equal(resourcesRead[1]) {
		return fmt.Errorf("expected data to be %v, got %v", resource.Data, resourcesRead[0])
	}

	return nil
}

func DestinationPluginTestSuiteRunner(t *testing.T, p *DestinationPlugin, spec specs.Destination, options ...DestinationPluginTestSuiteOption) {
	t.Helper()
	suite := &DestinationPluginTestSuite{}
	for _, option := range options {
		option(suite)
	}
	ctx := context.Background()
	logger := getTestLogger(t)

	t.Run("TestWriteOverwrite", func(t *testing.T) {
		t.Helper()
		if suite.skipTestOverwrite {
			t.Skip("skipping TestWriteOverwrite")
			return
		}
		if err := suite.destinationPluginTestWriteOverwrite(ctx, p, logger, spec); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TestWriteAppend", func(t *testing.T) {
		t.Helper()
		if suite.skipTestAppend {
			t.Skip("skipping TestWriteAppend")
			return
		}
		if err := suite.destinationPluginTestWriteAppend(ctx, p, logger, spec); err != nil {
			t.Fatal(err)
		}
	})
}
