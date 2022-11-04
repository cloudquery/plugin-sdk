package plugins

import "github.com/cloudquery/plugin-sdk/cqtypes"

type defaultTransformer struct {

}

func (p *defaultTransformer) TransformBool(v *cqtypes.Bool) interface{} {
	return v
}

func (p *defaultTransformer)  TransformBytea(v *cqtypes.Bytea) interface{} {
	return v
}

func (p *defaultTransformer)  TransformFloat8(v *cqtypes.Float8) interface{} {
	return v
}

func (p *defaultTransformer)  TransformInt8(v *cqtypes.Int8) interface{} {
	return v
}

func (p *defaultTransformer)  TransformInt8Array(v *cqtypes.Int8Array) interface{} {
	return v
}

func (p *defaultTransformer)  TransformJSON(v *cqtypes.JSON) interface{} {
	return v
}

func (p *defaultTransformer)  TransformText(v *cqtypes.Text) interface{} {
	return v
}

func (p *defaultTransformer)  TransformTextArray(v *cqtypes.TextArray) interface{} {
	return v
}

func (p *defaultTransformer)  TransformTimestamptz(v *cqtypes.Timestamptz) interface{} {
	return v
}

func (p *defaultTransformer)  TransformUUID(v *cqtypes.UUID) interface{} {
	return v
}

func (p *defaultTransformer)  TransformUUIDArray(v *cqtypes.UUIDArray) interface{} {
	return v
}

func (p *defaultTransformer)  TransformCIDR(v *cqtypes.CIDR) interface{} {
	return v
}

func (p *defaultTransformer)  TransformCIDRArray(v *cqtypes.CIDRArray) interface{} {
	return v
}

func (p *defaultTransformer)  TransformInet(v *cqtypes.Inet) interface{} {
	return v
}

func (p *defaultTransformer)  TransformInetArray(v *cqtypes.InetArray) interface{} {
	return v
}

func (p *defaultTransformer)  TransformMacaddr(v *cqtypes.Macaddr) interface{} {
	return v
}

func (p *defaultTransformer)  TransformMacaddrArray(v *cqtypes.MacaddrArray) interface{} {
	return v
}

