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

func TransformWithTransformer(transformer CQTypeTransformer, values CQTypes) []any {
	res := make([]any, len(values))
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

func (*DefaultTransformer) TransformBool(v *Bool) any {
	return v
}

func (*DefaultTransformer) TransformBytea(v *Bytea) any {
	return v
}

func (*DefaultTransformer) TransformFloat8(v *Float8) any {
	return v
}

func (*DefaultTransformer) TransformInt8(v *Int8) any {
	return v
}

func (*DefaultTransformer) TransformInt8Array(v *Int8Array) any {
	return v
}

func (*DefaultTransformer) TransformJSON(v *JSON) any {
	return v
}

func (*DefaultTransformer) TransformText(v *Text) any {
	return v
}

func (*DefaultTransformer) TransformTextArray(v *TextArray) any {
	return v
}

func (*DefaultTransformer) TransformTimestamptz(v *Timestamptz) any {
	return v
}

func (*DefaultTransformer) TransformUUID(v *UUID) any {
	return v
}

func (*DefaultTransformer) TransformUUIDArray(v *UUIDArray) any {
	return v
}

func (*DefaultTransformer) TransformCIDR(v *CIDR) any {
	return v
}

func (*DefaultTransformer) TransformCIDRArray(v *CIDRArray) any {
	return v
}

func (*DefaultTransformer) TransformInet(v *Inet) any {
	return v
}

func (*DefaultTransformer) TransformInetArray(v *InetArray) any {
	return v
}

func (*DefaultTransformer) TransformMacaddr(v *Macaddr) any {
	return v
}

func (*DefaultTransformer) TransformMacaddrArray(v *MacaddrArray) any {
	return v
}
