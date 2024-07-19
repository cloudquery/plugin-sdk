package reversertransformer

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"
)

// client is mostly used for testing the destination plugin.
type client struct {
	plugin.UnimplementedDestination
	plugin.UnimplementedSource
}

type Option func(*client)

type Spec struct {
}

func GetNewClient(options ...Option) plugin.NewClientFunc {
	c := &client{}
	for _, opt := range options {
		opt(c)
	}
	return func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return c, nil
	}
}

func (*client) GetSpec() any {
	return &Spec{}
}

func (*client) Close(context.Context) error {
	return nil
}

func (c *client) Transform(ctx context.Context, recvRecords <-chan arrow.Record, sendRecords chan<- arrow.Record) error {
	for {
		select {
		case record, ok := <-recvRecords:
			if !ok {
				return nil
			}
			reversedRecord, err := c.reverseStrings(record)
			if err != nil {
				return err
			}
			sendRecords <- reversedRecord
		case <-ctx.Done():
			return nil
		}
	}
}

func (*client) reverseStrings(record arrow.Record) (arrow.Record, error) {
	for i, column := range record.Columns() {
		if column.DataType().ID() != arrow.STRING {
			continue
		}
		newColumnData := []string{}
		for i := 0; i < column.Len(); i++ {
			if !column.IsValid(i) {
				continue
			}
			s := column.ValueStr(i)
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			newColumnData = append(newColumnData, string(runes))
		}
		fmt.Println("new column data is ", newColumnData)
		mem := memory.NewGoAllocator()
		bld := array.NewStringBuilder(mem)

		// create an array with 4 values, no null
		bld.AppendValues(newColumnData, nil)
		var err error
		record, err = record.SetColumn(i, bld.NewStringArray())
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}
