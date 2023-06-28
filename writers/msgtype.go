package writers

import (
	"reflect"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type msgType int

const (
	msgTypeUnset msgType = iota
	msgTypeMigrateTable
	msgTypeInsert
	msgTypeDeleteStale
)

func msgID(msg message.WriteMessage) msgType {
	switch msg.(type) {
	case *message.WriteMigrateTable:
		return msgTypeMigrateTable
	case *message.WriteInsert:
		return msgTypeInsert
	case *message.WriteDeleteStale:
		return msgTypeDeleteStale
	}
	panic("unknown message type: " + reflect.TypeOf(msg).Name())
}
