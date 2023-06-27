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

func msgID(msg message.Message) msgType {
	switch msg.(type) {
	case *message.MigrateTable:
		return msgTypeMigrateTable
	case *message.Insert:
		return msgTypeInsert
	case *message.DeleteStale:
		return msgTypeDeleteStale
	}
	panic("unknown message type: " + reflect.TypeOf(msg).Name())
}
