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

func TransformWithTransformer(transformer CQTypeTransformer, values CQTypes) ([]interface{}, error) {
	res := make([]interface{}, len(values))
	var err error
	for i, v := range values {
		switch v.Type() {
		case TypeBool:
			res[i], err = transformer.TransformBool(v.(*Bool))
		case TypeByteArray:
			res[i], err = transformer.TransformBytea(v.(*Bytea))
		case TypeCIDRArray:
			res[i], err = transformer.TransformCIDRArray(v.(*CIDRArray))
		case TypeCIDR:
			res[i], err = transformer.TransformCIDR(v.(*CIDR))
		case TypeFloat:
			res[i], err = transformer.TransformFloat8(v.(*Float8))
		case TypeInetArray:
			res[i], err = transformer.TransformInetArray(v.(*InetArray))
		case TypeInet:
			res[i], err = transformer.TransformInet(v.(*Inet))
		case TypeIntArray:
			res[i], err = transformer.TransformInt8Array(v.(*Int8Array))
		case TypeInt:
			res[i], err = transformer.TransformInt8(v.(*Int8))
		case TypeJSON:
			res[i], err = transformer.TransformJSON(v.(*JSON))
		case TypeMacAddrArray:
			res[i], err = transformer.TransformMacaddrArray(v.(*MacaddrArray))
		case TypeMacAddr:
			res[i], err = transformer.TransformMacaddr(v.(*Macaddr))
		case TypeStringArray:
			res[i], err = transformer.TransformTextArray(v.(*TextArray))
		case TypeString:
			res[i], err = transformer.TransformText(v.(*Text))
		case TypeTimestamp:
			res[i], err = transformer.TransformTimestamptz(v.(*Timestamptz))
		case TypeUUIDArray:
			res[i], err = transformer.TransformUUIDArray(v.(*UUIDArray))
		case TypeUUID:
			res[i], err = transformer.TransformUUID(v.(*UUID))
		default:
			panic("unknown type " + v.Type().String())
		}

		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

type DefaultTransformer struct {
}

func (*DefaultTransformer) TransformBool(v *Bool) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformBytea(v *Bytea) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformFloat8(v *Float8) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformInt8(v *Int8) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformInt8Array(v *Int8Array) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformJSON(v *JSON) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformText(v *Text) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformTextArray(v *TextArray) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformTimestamptz(v *Timestamptz) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformUUID(v *UUID) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformUUIDArray(v *UUIDArray) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformCIDR(v *CIDR) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformCIDRArray(v *CIDRArray) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformInet(v *Inet) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformInetArray(v *InetArray) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformMacaddr(v *Macaddr) (interface{}, error) {
	return v, nil
}

func (*DefaultTransformer) TransformMacaddrArray(v *MacaddrArray) (interface{}, error) {
	return v, nil
}
