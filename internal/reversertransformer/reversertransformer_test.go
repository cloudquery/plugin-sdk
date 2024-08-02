package reversertransformer

import (
	"context"
	"io"
	"testing"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/apache/arrow/go/v17/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	internalPlugin "github.com/cloudquery/plugin-sdk/v4/internal/servers/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestReverserTransformer(t *testing.T) {
	p := plugin.NewPlugin("test", "development", GetNewClient())
	s := internalPlugin.Server{
		Plugin: p,
	}
	_, err := s.Init(context.Background(), &pb.Init_Request{
		Spec:         []byte("{}"),
		NoConnection: true,
		InvocationId: "26b550f9-c6f8-4b4b-9ec4-773bab288ee6",
	})
	require.NoError(t, err)
	requests := makeRequestsFromStrings("hello", "world")
	stream := mockTransformServer{incomingMessages: requests}
	require.NoError(t, s.Transform(&stream))
	require.Equal(t, 2, len(stream.outgoingMessages))

	record1, err := pb.NewRecordFromBytes(stream.outgoingMessages[0].Record)
	require.NoError(t, err)
	record2, err := pb.NewRecordFromBytes(stream.outgoingMessages[1].Record)
	require.NoError(t, err)

	require.Equal(t, "olleh", record1.Column(0).ValueStr(0))
	require.Equal(t, "dlrow", record2.Column(0).ValueStr(0))
}

func makeRequestsFromStrings(s ...string) []*pb.Transform_Request {
	requests := make([]*pb.Transform_Request, len(s))
	for i, str := range s {
		requests[i] = makeRequestFromString(str)
	}
	return requests
}

func makeRequestFromString(s string) *pb.Transform_Request {
	record := makeRecordFromString(s)
	bs, _ := pb.RecordToBytes(record)
	return &pb.Transform_Request{Record: bs}
}

func makeRecordFromString(s string) arrow.Record {
	str := array.NewStringBuilder(memory.DefaultAllocator)
	str.AppendString(s)
	arr := str.NewStringArray()
	schema := arrow.NewSchema([]arrow.Field{{Name: "col1", Type: arrow.BinaryTypes.String}}, nil)

	return array.NewRecord(schema, []arrow.Array{arr}, 1)
}

type mockTransformServer struct {
	grpc.ServerStream
	incomingMessages []*pb.Transform_Request
	outgoingMessages []*pb.Transform_Response
}

func (*mockTransformServer) SendAndClose(*pb.Transform_Response) error {
	return nil
}
func (s *mockTransformServer) Recv() (*pb.Transform_Request, error) {
	if len(s.incomingMessages) > 0 {
		msg := s.incomingMessages[0]
		s.incomingMessages = s.incomingMessages[1:]
		return msg, nil
	}
	return nil, io.EOF
}
func (s *mockTransformServer) Send(resp *pb.Transform_Response) error {
	s.outgoingMessages = append(s.outgoingMessages, resp)
	return nil
}
func (*mockTransformServer) SetHeader(metadata.MD) error {
	return nil
}
func (*mockTransformServer) SendHeader(metadata.MD) error {
	return nil
}
func (*mockTransformServer) SetTrailer(metadata.MD) {
}
func (mockTransformServer) Context() context.Context {
	return context.Background()
}
func (mockTransformServer) SendMsg(any) error {
	return nil
}
func (mockTransformServer) RecvMsg(any) error {
	return nil
}
