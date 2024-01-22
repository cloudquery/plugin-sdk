package scalar

import (
	"net"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type Mac struct {
	Valid bool
	Value net.HardwareAddr
}

func (s *Mac) IsValid() bool {
	return s.Valid
}

func (*Mac) DataType() arrow.DataType {
	return types.ExtensionTypes.MAC
}

func (s *Mac) String() string {
	if !s.Valid {
		return nullValueStr
	}
	return s.Value.String()
}

func (s *Mac) Equal(rhs Scalar) bool {
	if rhs == nil {
		return false
	}
	r, ok := rhs.(*Mac)
	if !ok {
		return false
	}
	return s.Valid == r.Valid && s.Value.String() == r.Value.String()
}

func (s *Mac) Get() any {
	return s.Value
}

func (s *Mac) Set(val any) error {
	if val == nil {
		s.Valid = false
		return nil
	}

	if sc, ok := val.(Scalar); ok {
		if !sc.IsValid() {
			s.Valid = false
			return nil
		}
		return s.Set(sc.Get())
	}

	switch value := val.(type) {
	case net.HardwareAddr:
		addr := make(net.HardwareAddr, len(value))
		copy(addr, value)
		s.Value = addr
	case string:
		addr, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		s.Value = addr
	case *net.HardwareAddr:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	case *string:
		if value == nil {
			s.Valid = false
			return nil
		}
		return s.Set(*value)
	default:
		if originalSrc, ok := underlyingPtrType(value); ok {
			return s.Set(originalSrc)
		}
		return &ValidationError{Type: types.ExtensionTypes.MAC, Msg: noConversion, Value: value}
	}

	s.Valid = true
	return nil
}
