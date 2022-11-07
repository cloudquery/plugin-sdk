package schema

type CQTypeTransformer interface {
	BoolTransformer
	ByteaTransformer
	CIDRArrayTransformer
	CIDRTransformer
	Float8Transformer
	InetArrayTransformer
	InetTransformer
	Int8ArrayTransformer
	Int8Transformer
	JSONTransformer
	MacaddrArrayTransformer
	MacaddrTransformer
	TextArrayTransformer
	TextTransformer
	TimestamptzTransformer
	UUIDArrayTransformer
	UUIDTransformer
}

func TransformWithTransformer(transformer CQTypeTransformer, values CQTypes) []interface{} {
	res := make([]interface{}, len(values))
	for i, v := range values {
		switch v.Type() {
		case TypeBool:
			res[i] = transformer.TransformBool(v.(*Bool))
		case TypeByteArray:
			res[i] = transformer.TransformBytea(v.(*Bytea))
		case TypeCIDRArray:
			res[i] = transformer.TransformCIDRArray(v.(*CIDRArray))
		case TypeCIDR:
			res[i] = transformer.TransformCIDR(v.(*CIDR))
		case TypeFloat:
			res[i] = transformer.TransformFloat8(v.(*Float8))
		case TypeInetArray:
			res[i] = transformer.TransformInetArray(v.(*InetArray))
		case TypeInet:
			res[i] = transformer.TransformInet(v.(*Inet))
		case TypeIntArray:
			res[i] = transformer.TransformInt8Array(v.(*Int8Array))
		case TypeInt:
			res[i] = transformer.TransformInt8(v.(*Int8))
		case TypeJSON:
			res[i] = transformer.TransformJSON(v.(*JSON))
		case TypeMacAddrArray:
			res[i] = transformer.TransformMacaddrArray(v.(*MacaddrArray))
		case TypeMacAddr:
			res[i] = transformer.TransformMacaddr(v.(*Macaddr))
		case TypeStringArray:
			res[i] = transformer.TransformTextArray(v.(*TextArray))
		case TypeString:
			res[i] = transformer.TransformText(v.(*Text))
		case TypeTimestamp:
			res[i] = transformer.TransformTimestamptz(v.(*Timestamptz))
		case TypeUUIDArray:
			res[i] = transformer.TransformUUIDArray(v.(*UUIDArray))
		case TypeUUID:
			res[i] = transformer.TransformUUID(v.(*UUID))
		default:
			panic("unknown type " + v.Type().String())
		}
	}
	return res
}

type DefaultTransformer struct {
}

func (*DefaultTransformer) TransformBool(v *Bool) interface{} {
	return v
}

func (*DefaultTransformer) TransformBytea(v *Bytea) interface{} {
	return v
}

func (*DefaultTransformer) TransformFloat8(v *Float8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8(v *Int8) interface{} {
	return v
}

func (*DefaultTransformer) TransformInt8Array(v *Int8Array) interface{} {
	return v
}

func (*DefaultTransformer) TransformJSON(v *JSON) interface{} {
	return v
}

func (*DefaultTransformer) TransformText(v *Text) interface{} {
	return v
}

func (*DefaultTransformer) TransformTextArray(v *TextArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformTimestamptz(v *Timestamptz) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUID(v *UUID) interface{} {
	return v
}

func (*DefaultTransformer) TransformUUIDArray(v *UUIDArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDR(v *CIDR) interface{} {
	return v
}

func (*DefaultTransformer) TransformCIDRArray(v *CIDRArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformInet(v *Inet) interface{} {
	return v
}

func (*DefaultTransformer) TransformInetArray(v *InetArray) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddr(v *Macaddr) interface{} {
	return v
}

func (*DefaultTransformer) TransformMacaddrArray(v *MacaddrArray) interface{} {
	return v
}
