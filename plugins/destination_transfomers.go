package plugins

import "github.com/cloudquery/plugin-sdk/v1/schema"

type DefaultTransformer struct {
}

func (*DefaultTransformer) TransformBool(v *schema.Bool) interface{} {
	return v
}

func (*DefaultTransformer) TransformBytea(v *schema.Bytea) interface{} {
	return v
}

func (*DefaultTransformer) TransformFloat8(v *schema.Float8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8(v *schema.Int8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8Array(v *schema.Int8Array) interface{} {
	return v
}

func (*DefaultTransformer) TransformJSON(v *schema.JSON) interface{} {
	return v
}

func (*DefaultTransformer) TransformText(v *schema.Text) interface{} {
	return v
}

func (*DefaultTransformer) TransformTextArray(v *schema.TextArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformTimestamptz(v *schema.Timestamptz) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUID(v *schema.UUID) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUIDArray(v *schema.UUIDArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDR(v *schema.CIDR) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDRArray(v *schema.CIDRArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformInet(v *schema.Inet) interface{} {
	return v
}

func (*DefaultTransformer) TransformInetArray(v *schema.InetArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddr(v *schema.Macaddr) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddrArray(v *schema.MacaddrArray) interface{} {
	return v
}
