package cqtypes

type CQTypeTransformer interface {
	BoolTransformer
	ByteaTransformer
	CIDRArrayTransformer
	CIDRTransformer
	Float8Transformer
	InetArrayTransformer
	InetTransformer
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
