package plugins

import "github.com/cloudquery/plugin-sdk/cqtypes"

type DefaultTransformer struct {
}

func (*DefaultTransformer) TransformBool(v *cqtypes.Bool) interface{} {
	return v
}

func (*DefaultTransformer) TransformBytea(v *cqtypes.Bytea) interface{} {
	return v
}

func (*DefaultTransformer) TransformFloat8(v *cqtypes.Float8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8(v *cqtypes.Int8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8Array(v *cqtypes.Int8Array) interface{} {
	return v
}

func (*DefaultTransformer) TransformJSON(v *cqtypes.JSON) interface{} {
	return v
}

func (*DefaultTransformer) TransformText(v *cqtypes.Text) interface{} {
	return v
}

func (*DefaultTransformer) TransformTextArray(v *cqtypes.TextArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformTimestamptz(v *cqtypes.Timestamptz) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUID(v *cqtypes.UUID) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUIDArray(v *cqtypes.UUIDArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDR(v *cqtypes.CIDR) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDRArray(v *cqtypes.CIDRArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformInet(v *cqtypes.Inet) interface{} {
	return v
}

func (*DefaultTransformer) TransformInetArray(v *cqtypes.InetArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddr(v *cqtypes.Macaddr) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddrArray(v *cqtypes.MacaddrArray) interface{} {
	return v
}
