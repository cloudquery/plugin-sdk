package writers

import (
	"reflect"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type MsgType int

const (
	MsgTypeUnset MsgType = iota
	MsgTypeMigrateTable
	MsgTypeInsert
	MsgTypeDeleteStale
	MsgTypeDeleteRecord
)

func MsgID(msg message.WriteMessage) MsgType {
	switch msg.(type) {
	case *message.WriteMigrateTable:
		return MsgTypeMigrateTable
	case *message.WriteInsert:
		return MsgTypeInsert
	case *message.WriteDeleteStale:
		return MsgTypeDeleteStale
	case *message.WriteDeleteRecord:
		return MsgTypeDeleteRecord
	}
	panic("unknown message type: " + reflect.TypeOf(msg).Name())
}
