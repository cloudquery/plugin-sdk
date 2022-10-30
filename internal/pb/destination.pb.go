// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: internal/pb/destination.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Migrate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Migrate) Reset() {
	*x = Migrate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate) ProtoMessage() {}

func (x *Migrate) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate.ProtoReflect.Descriptor instead.
func (*Migrate) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0}
}

type Write struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Write) Reset() {
	*x = Write{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write) ProtoMessage() {}

func (x *Write) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write.ProtoReflect.Descriptor instead.
func (*Write) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1}
}

type Close struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Close) Reset() {
	*x = Close{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Close) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close) ProtoMessage() {}

func (x *Close) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close.ProtoReflect.Descriptor instead.
func (*Close) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{2}
}

type DeleteStale struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteStale) Reset() {
	*x = DeleteStale{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteStale) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale) ProtoMessage() {}

func (x *DeleteStale) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale.ProtoReflect.Descriptor instead.
func (*DeleteStale) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{3}
}

type GetDestinationMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDestinationMetrics) Reset() {
	*x = GetDestinationMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDestinationMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics) ProtoMessage() {}

func (x *GetDestinationMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{4}
}

type Migrate_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Tables  []byte `protobuf:"bytes,3,opt,name=tables,proto3" json:"tables,omitempty"`
}

func (x *Migrate_Request) Reset() {
	*x = Migrate_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Request) ProtoMessage() {}

func (x *Migrate_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Request.ProtoReflect.Descriptor instead.
func (*Migrate_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Migrate_Request) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Migrate_Request) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Migrate_Request) GetTables() []byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type Migrate_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Migrate_Response) Reset() {
	*x = Migrate_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Response) ProtoMessage() {}

func (x *Migrate_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Response.ProtoReflect.Descriptor instead.
func (*Migrate_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0, 1}
}

type Write_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled *schema.Resources
	Resource  []byte                 `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	Source    string                 `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Write_Request) Reset() {
	*x = Write_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Request) ProtoMessage() {}

func (x *Write_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Request.ProtoReflect.Descriptor instead.
func (*Write_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Write_Request) GetResource() []byte {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *Write_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Write_Request) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type Write_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FailedWrites uint64 `protobuf:"varint,1,opt,name=failed_writes,json=failedWrites,proto3" json:"failed_writes,omitempty"`
}

func (x *Write_Response) Reset() {
	*x = Write_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Response) ProtoMessage() {}

func (x *Write_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Response.ProtoReflect.Descriptor instead.
func (*Write_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Write_Response) GetFailedWrites() uint64 {
	if x != nil {
		return x.FailedWrites
	}
	return 0
}

type Close_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Close_Request) Reset() {
	*x = Close_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Close_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close_Request) ProtoMessage() {}

func (x *Close_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close_Request.ProtoReflect.Descriptor instead.
func (*Close_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{2, 0}
}

type Close_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Close_Response) Reset() {
	*x = Close_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Close_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close_Response) ProtoMessage() {}

func (x *Close_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close_Response.ProtoReflect.Descriptor instead.
func (*Close_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{2, 1}
}

type DeleteStale_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source    string                 `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Tables    []byte                 `protobuf:"bytes,3,opt,name=tables,proto3" json:"tables,omitempty"`
}

func (x *DeleteStale_Request) Reset() {
	*x = DeleteStale_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteStale_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale_Request) ProtoMessage() {}

func (x *DeleteStale_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale_Request.ProtoReflect.Descriptor instead.
func (*DeleteStale_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{3, 0}
}

func (x *DeleteStale_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *DeleteStale_Request) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *DeleteStale_Request) GetTables() []byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type DeleteStale_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FailedDeletes uint64 `protobuf:"varint,1,opt,name=failed_deletes,json=failedDeletes,proto3" json:"failed_deletes,omitempty"`
}

func (x *DeleteStale_Response) Reset() {
	*x = DeleteStale_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteStale_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale_Response) ProtoMessage() {}

func (x *DeleteStale_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale_Response.ProtoReflect.Descriptor instead.
func (*DeleteStale_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{3, 1}
}

func (x *DeleteStale_Response) GetFailedDeletes() uint64 {
	if x != nil {
		return x.FailedDeletes
	}
	return 0
}

type GetDestinationMetrics_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDestinationMetrics_Request) Reset() {
	*x = GetDestinationMetrics_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDestinationMetrics_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics_Request) ProtoMessage() {}

func (x *GetDestinationMetrics_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics_Request.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{4, 0}
}

type GetDestinationMetrics_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled json of plugins.DestinationMetrics
	Metrics []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *GetDestinationMetrics_Response) Reset() {
	*x = GetDestinationMetrics_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDestinationMetrics_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics_Response) ProtoMessage() {}

func (x *GetDestinationMetrics_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics_Response.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{4, 1}
}

func (x *GetDestinationMetrics_Response) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_internal_pb_destination_proto protoreflect.FileDescriptor

var file_internal_pb_destination_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x62, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x66, 0x0a, 0x07, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x1a, 0x4f, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x0a, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xb1, 0x01, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74,
	0x65, 0x1a, 0x77, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x2f, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x66,
	0x61, 0x69, 0x6c, 0x65, 0x64, 0x57, 0x72, 0x69, 0x74, 0x65, 0x73, 0x22, 0x1e, 0x0a, 0x05, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xb5, 0x01, 0x0a, 0x0b,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x6c, 0x65, 0x1a, 0x73, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x38,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
	0x1a, 0x31, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e,
	0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x73, 0x22, 0x48, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x09, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x32, 0xfa, 0x04,
	0x0a, 0x0b, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5b, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a,
	0x07, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28,
	0x01, 0x12, 0x34, 0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x59, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x24, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x44,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pb_destination_proto_rawDescOnce sync.Once
	file_internal_pb_destination_proto_rawDescData = file_internal_pb_destination_proto_rawDesc
)

func file_internal_pb_destination_proto_rawDescGZIP() []byte {
	file_internal_pb_destination_proto_rawDescOnce.Do(func() {
		file_internal_pb_destination_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pb_destination_proto_rawDescData)
	})
	return file_internal_pb_destination_proto_rawDescData
}

var file_internal_pb_destination_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_internal_pb_destination_proto_goTypes = []interface{}{
	(*Migrate)(nil),                        // 0: proto.Migrate
	(*Write)(nil),                          // 1: proto.Write
	(*Close)(nil),                          // 2: proto.Close
	(*DeleteStale)(nil),                    // 3: proto.DeleteStale
	(*GetDestinationMetrics)(nil),          // 4: proto.GetDestinationMetrics
	(*Migrate_Request)(nil),                // 5: proto.Migrate.Request
	(*Migrate_Response)(nil),               // 6: proto.Migrate.Response
	(*Write_Request)(nil),                  // 7: proto.Write.Request
	(*Write_Response)(nil),                 // 8: proto.Write.Response
	(*Close_Request)(nil),                  // 9: proto.Close.Request
	(*Close_Response)(nil),                 // 10: proto.Close.Response
	(*DeleteStale_Request)(nil),            // 11: proto.DeleteStale.Request
	(*DeleteStale_Response)(nil),           // 12: proto.DeleteStale.Response
	(*GetDestinationMetrics_Request)(nil),  // 13: proto.GetDestinationMetrics.Request
	(*GetDestinationMetrics_Response)(nil), // 14: proto.GetDestinationMetrics.Response
	(*timestamppb.Timestamp)(nil),          // 15: google.protobuf.Timestamp
	(*GetProtocolVersion_Request)(nil),     // 16: proto.GetProtocolVersion.Request
	(*GetName_Request)(nil),                // 17: proto.GetName.Request
	(*GetVersion_Request)(nil),             // 18: proto.GetVersion.Request
	(*Configure_Request)(nil),              // 19: proto.Configure.Request
	(*GetProtocolVersion_Response)(nil),    // 20: proto.GetProtocolVersion.Response
	(*GetName_Response)(nil),               // 21: proto.GetName.Response
	(*GetVersion_Response)(nil),            // 22: proto.GetVersion.Response
	(*Configure_Response)(nil),             // 23: proto.Configure.Response
}
var file_internal_pb_destination_proto_depIdxs = []int32{
	15, // 0: proto.Write.Request.timestamp:type_name -> google.protobuf.Timestamp
	15, // 1: proto.DeleteStale.Request.timestamp:type_name -> google.protobuf.Timestamp
	16, // 2: proto.Destination.GetProtocolVersion:input_type -> proto.GetProtocolVersion.Request
	17, // 3: proto.Destination.GetName:input_type -> proto.GetName.Request
	18, // 4: proto.Destination.GetVersion:input_type -> proto.GetVersion.Request
	19, // 5: proto.Destination.Configure:input_type -> proto.Configure.Request
	5,  // 6: proto.Destination.Migrate:input_type -> proto.Migrate.Request
	7,  // 7: proto.Destination.Write:input_type -> proto.Write.Request
	9,  // 8: proto.Destination.Close:input_type -> proto.Close.Request
	11, // 9: proto.Destination.DeleteStale:input_type -> proto.DeleteStale.Request
	13, // 10: proto.Destination.GetMetrics:input_type -> proto.GetDestinationMetrics.Request
	20, // 11: proto.Destination.GetProtocolVersion:output_type -> proto.GetProtocolVersion.Response
	21, // 12: proto.Destination.GetName:output_type -> proto.GetName.Response
	22, // 13: proto.Destination.GetVersion:output_type -> proto.GetVersion.Response
	23, // 14: proto.Destination.Configure:output_type -> proto.Configure.Response
	6,  // 15: proto.Destination.Migrate:output_type -> proto.Migrate.Response
	8,  // 16: proto.Destination.Write:output_type -> proto.Write.Response
	10, // 17: proto.Destination.Close:output_type -> proto.Close.Response
	12, // 18: proto.Destination.DeleteStale:output_type -> proto.DeleteStale.Response
	14, // 19: proto.Destination.GetMetrics:output_type -> proto.GetDestinationMetrics.Response
	11, // [11:20] is the sub-list for method output_type
	2,  // [2:11] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_internal_pb_destination_proto_init() }
func file_internal_pb_destination_proto_init() {
	if File_internal_pb_destination_proto != nil {
		return
	}
	file_internal_pb_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_internal_pb_destination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Close); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteStale); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDestinationMetrics); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Close_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Close_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteStale_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteStale_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDestinationMetrics_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_pb_destination_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDestinationMetrics_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_pb_destination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pb_destination_proto_goTypes,
		DependencyIndexes: file_internal_pb_destination_proto_depIdxs,
		MessageInfos:      file_internal_pb_destination_proto_msgTypes,
	}.Build()
	File_internal_pb_destination_proto = out.File
	file_internal_pb_destination_proto_rawDesc = nil
	file_internal_pb_destination_proto_goTypes = nil
	file_internal_pb_destination_proto_depIdxs = nil
}
