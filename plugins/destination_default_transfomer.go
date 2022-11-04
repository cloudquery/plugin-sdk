package plugins

import "github.com/cloudquery/plugin-sdk/cqtypes"

type defaultTransformer struct {
}

func (*defaultTransformer) TransformBool(v *cqtypes.Bool) interface{} {
	return v
}

func (*defaultTransformer) TransformBytea(v *cqtypes.Bytea) interface{} {
	return v
}

func (*defaultTransformer) TransformFloat8(v *cqtypes.Float8) interface{} {
	return v
}

func (*defaultTransformer) TransformInt8(v *cqtypes.Int8) interface{} {
	return v
}

func (*defaultTransformer) TransformInt8Array(v *cqtypes.Int8Array) interface{} {
	return v
}

func (*defaultTransformer) TransformJSON(v *cqtypes.JSON) interface{} {
	return v
}

func (*defaultTransformer) TransformText(v *cqtypes.Text) interface{} {
	return v
}

func (*defaultTransformer) TransformTextArray(v *cqtypes.TextArray) interface{} {
	return v
}

func (*defaultTransformer) TransformTimestamptz(v *cqtypes.Timestamptz) interface{} {
	return v
}

func (*defaultTransformer) TransformUUID(v *cqtypes.UUID) interface{} {
	return v
}

func (*defaultTransformer) TransformUUIDArray(v *cqtypes.UUIDArray) interface{} {
	return v
}

func (*defaultTransformer) TransformCIDR(v *cqtypes.CIDR) interface{} {
	return v
}

func (*defaultTransformer) TransformCIDRArray(v *cqtypes.CIDRArray) interface{} {
	return v
}

func (*defaultTransformer) TransformInet(v *cqtypes.Inet) interface{} {
	return v
}

func (*defaultTransformer) TransformInetArray(v *cqtypes.InetArray) interface{} {
	return v
}

func (*defaultTransformer) TransformMacaddr(v *cqtypes.Macaddr) interface{} {
	return v
}

func (*defaultTransformer) TransformMacaddrArray(v *cqtypes.MacaddrArray) interface{} {
	return v
}
