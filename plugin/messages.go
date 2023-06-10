package plugin

import (
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type MessageType int

const (
	// Create table
	MessageTypeCreate MessageType = iota
	// Insert record
	MessageTypeInsert
	// Insert or update record
	MessageTypeUpsert
	// Delete rows
	MessageTypeDelete
)

type MessageCreateTable struct {
	Table *schema.Table
	Force bool
}

func (*MessageCreateTable) Type() MessageType {
	return MessageTypeCreate
}

type MessageInsert struct {
	Record  arrow.Record
	Columns []string
	Upsert  bool
}

func (*MessageInsert) Type() MessageType {
	return MessageTypeInsert
}

type Operator int

const (
	OperatorEqual Operator = iota
	OperatorNotEqual
	OperatorGreaterThan
	OperatorGreaterThanOrEqual
	OperatorLessThan
	OperatorLessThanOrEqual
)

type WhereClause struct {
	Column   string
	Operator Operator
	Value    string
}

type MessageDelete struct {
	Record arrow.Record
	// currently delete only supports and where clause as we don't support
	// full AST parsing
	WhereClauses []WhereClause
}

func (*MessageDelete) Type() MessageType {
	return MessageTypeDelete
}

type Message interface {
	Type() MessageType
}

type Messages []Message

func (m Messages) InsertItems() int64 {
	items := int64(0)
	for _, msg := range m {
		switch msg.Type() {
		case MessageTypeInsert:
			msgInsert := msg.(*MessageInsert)
			items += msgInsert.Record.NumRows()
		}
	}
	return items
}
