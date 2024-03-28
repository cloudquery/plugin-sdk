package plugin

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func (s *WriterTestSuite) testUpsertBasic(ctx context.Context) error {
	tableName := s.tableNameForTest("upsert_basic")
	table := &schema.Table{
		Name: tableName,
		Columns: []schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, NotNull: true},
			{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true},
		},
	}
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	bldr.Field(0).(*array.Int64Builder).Append(1)
	bldr.Field(1).(*array.StringBuilder).Append("foo")
	record := bldr.NewRecord()

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	record = s.handleNulls(record) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	totalItems := TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}

	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: record,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	records, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	totalItems = TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected 1 item, got %d", totalItems)
	}
	if diff := RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{record}); diff != "" {
		return fmt.Errorf("record differs: %s", diff)
	}
	return nil
}

func (s *WriterTestSuite) testUpsertAll(ctx context.Context) error {
	const rowsPerRecord = 10
	tableName := s.tableNameForTest("upsert_all")
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns = append(table.Columns, schema.Column{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKey: true})
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	tg := schema.NewTestDataGenerator(0)
	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:            rowsPerRecord,
		TimePrecision:      s.genDatOptions.TimePrecision,
		UseHomogeneousType: s.useHomogeneousTypes,
	})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: normalRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	normalRecord = s.handleNulls(normalRecord) // we process nulls after writing

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	sortRecords(table, records, "id")

	totalItems := TotalRows(records)
	if totalItems != rowsPerRecord {
		return fmt.Errorf("expected items after initial insert: %d, got %d", rowsPerRecord, totalItems)
	}

	if diff := RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{normalRecord}); diff != "" {
		return fmt.Errorf("record differs after insert: %s", diff)
	}

	tg.Reset()
	nullRecord := tg.Generate(table, schema.GenTestDataOptions{MaxRows: rowsPerRecord, TimePrecision: s.genDatOptions.TimePrecision, NullRows: true})
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: nullRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	nullRecord = s.handleNulls(nullRecord) // we process nulls after writing

	records, err = s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	sortRecords(table, records, "id")

	totalItems = TotalRows(records)
	if totalItems != rowsPerRecord {
		return fmt.Errorf("expected items after upsert: %d, got %d", rowsPerRecord, totalItems)
	}

	if diff := RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{nullRecord}); diff != "" {
		return fmt.Errorf("record differs after upsert (columns should be null): %s", diff)
	}

	return nil
}

func (s *WriterTestSuite) testInsertDuplicatePK(ctx context.Context) error {
	const rowsPerRecord = 10
	tableName := s.tableNameForTest("upsert_duplicate_pk")
	table := schema.TestTable(tableName, s.genDatOptions)
	table.Columns.Get("id").PrimaryKey = true
	if err := s.plugin.writeOne(ctx, &message.WriteMigrateTable{
		Table: table,
	}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	tg := schema.NewTestDataGenerator(0)

	normalRecord := tg.Generate(table, schema.GenTestDataOptions{
		MaxRows:            rowsPerRecord,
		TimePrecision:      s.genDatOptions.TimePrecision,
		UseHomogeneousType: s.useHomogeneousTypes,
		DuplicateData:      true,
	})

	// normalRecord
	if err := s.plugin.writeOne(ctx, &message.WriteInsert{
		Record: normalRecord,
	}); err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	records, err := s.plugin.readAll(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to readAll: %w", err)
	}
	sortRecords(table, records, "id")

	totalItems := TotalRows(records)
	if totalItems != 1 {
		return fmt.Errorf("expected items after initial insert: %d, got %d", 1, totalItems)
	}

	if diff := RecordsDiff(table.ToArrowSchema(), records, []arrow.Record{extractLastRowFromRecord(table, normalRecord)}); diff != "" {
		return fmt.Errorf("record differs after insert: %s", diff)
	}

	return nil
}

func extractLastRowFromRecord(table *schema.Table, existingRecord arrow.Record) arrow.Record {
	sc := table.ToArrowSchema()
	var lastRecord []arrow.Record
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	for i, c := range table.Columns {
		col := existingRecord.Column(i)
		err := bldr.Field(i).AppendValueFromString(col.ValueStr(int(existingRecord.NumRows()) - 1))
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal json `%v` for column %v: %v", col.ValueStr(int(existingRecord.NumRows())-1), c.Name, err))
		}
	}
	lastRecord = append(lastRecord, bldr.NewRecord())
	bldr.Release()

	arrowTable := array.NewTableFromRecords(sc, lastRecord)
	columns := make([]arrow.Array, sc.NumFields())
	for n := 0; n < sc.NumFields(); n++ {
		concatenated, err := array.Concatenate(arrowTable.Column(n).Data().Chunks(), memory.DefaultAllocator)
		if err != nil {
			panic(fmt.Sprintf("failed to concatenate arrays: %v", err))
		}
		columns[n] = concatenated
	}

	return array.NewRecord(sc, columns, -1)
}
