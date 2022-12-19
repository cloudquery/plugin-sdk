// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: internal/pb/source.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Sync struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Sync) Reset() {
	*x = Sync{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync) ProtoMessage() {}

func (x *Sync) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync.ProtoReflect.Descriptor instead.
func (*Sync) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{0}
}

type Sync2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Sync2) Reset() {
	*x = Sync2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync2) ProtoMessage() {}

func (x *Sync2) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync2.ProtoReflect.Descriptor instead.
func (*Sync2) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{1}
}

type GetSyncSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSyncSummary) Reset() {
	*x = GetSyncSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSyncSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSyncSummary) ProtoMessage() {}

func (x *GetSyncSummary) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSyncSummary.ProtoReflect.Descriptor instead.
func (*GetSyncSummary) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{2}
}

type GetTables struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTables) Reset() {
	*x = GetTables{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTables) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTables) ProtoMessage() {}

func (x *GetTables) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTables.ProtoReflect.Descriptor instead.
func (*GetTables) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{3}
}

type GetTablesForSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTablesForSpec) Reset() {
	*x = GetTablesForSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTablesForSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTablesForSpec) ProtoMessage() {}

func (x *GetTablesForSpec) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTablesForSpec.ProtoReflect.Descriptor instead.
func (*GetTablesForSpec) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{4}
}

type GetSourceMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSourceMetrics) Reset() {
	*x = GetSourceMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSourceMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSourceMetrics) ProtoMessage() {}

func (x *GetSourceMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSourceMetrics.ProtoReflect.Descriptor instead.
func (*GetSourceMetrics) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{5}
}

type Sync_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Spec []byte `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *Sync_Request) Reset() {
	*x = Sync_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync_Request) ProtoMessage() {}

func (x *Sync_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync_Request.ProtoReflect.Descriptor instead.
func (*Sync_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Sync_Request) GetSpec() []byte {
	if x != nil {
		return x.Spec
	}
	return nil
}

type Sync_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled *schema.Resources
	Resource []byte `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *Sync_Response) Reset() {
	*x = Sync_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync_Response) ProtoMessage() {}

func (x *Sync_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync_Response.ProtoReflect.Descriptor instead.
func (*Sync_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Sync_Response) GetResource() []byte {
	if x != nil {
		return x.Resource
	}
	return nil
}

type Sync2_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Spec []byte `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *Sync2_Request) Reset() {
	*x = Sync2_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync2_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync2_Request) ProtoMessage() {}

func (x *Sync2_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync2_Request.ProtoReflect.Descriptor instead.
func (*Sync2_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Sync2_Request) GetSpec() []byte {
	if x != nil {
		return x.Spec
	}
	return nil
}

type Sync2_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled *schema.Resources
	Resource []byte `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *Sync2_Response) Reset() {
	*x = Sync2_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sync2_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sync2_Response) ProtoMessage() {}

func (x *Sync2_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sync2_Response.ProtoReflect.Descriptor instead.
func (*Sync2_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Sync2_Response) GetResource() []byte {
	if x != nil {
		return x.Resource
	}
	return nil
}

type GetSyncSummary_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSyncSummary_Request) Reset() {
	*x = GetSyncSummary_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSyncSummary_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSyncSummary_Request) ProtoMessage() {}

func (x *GetSyncSummary_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSyncSummary_Request.ProtoReflect.Descriptor instead.
func (*GetSyncSummary_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{2, 0}
}

type GetSyncSummary_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled *schema.SyncSummary
	Summary []byte `protobuf:"bytes,1,opt,name=summary,proto3" json:"summary,omitempty"`
}

func (x *GetSyncSummary_Response) Reset() {
	*x = GetSyncSummary_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSyncSummary_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSyncSummary_Response) ProtoMessage() {}

func (x *GetSyncSummary_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSyncSummary_Response.ProtoReflect.Descriptor instead.
func (*GetSyncSummary_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{2, 1}
}

func (x *GetSyncSummary_Response) GetSummary() []byte {
	if x != nil {
		return x.Summary
	}
	return nil
}

type GetTables_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTables_Request) Reset() {
	*x = GetTables_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTables_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTables_Request) ProtoMessage() {}

func (x *GetTables_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTables_Request.ProtoReflect.Descriptor instead.
func (*GetTables_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{3, 0}
}

type GetTables_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Marshalled []*schema.Table
	Tables []byte `protobuf:"bytes,3,opt,name=tables,proto3" json:"tables,omitempty"`
}

func (x *GetTables_Response) Reset() {
	*x = GetTables_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTables_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTables_Response) ProtoMessage() {}

func (x *GetTables_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTables_Response.ProtoReflect.Descriptor instead.
func (*GetTables_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{3, 1}
}

func (x *GetTables_Response) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetTables_Response) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *GetTables_Response) GetTables() []byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type GetTablesForSpec_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Marshalled specs.Source
	Spec []byte `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *GetTablesForSpec_Request) Reset() {
	*x = GetTablesForSpec_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTablesForSpec_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTablesForSpec_Request) ProtoMessage() {}

func (x *GetTablesForSpec_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTablesForSpec_Request.ProtoReflect.Descriptor instead.
func (*GetTablesForSpec_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{4, 0}
}

func (x *GetTablesForSpec_Request) GetSpec() []byte {
	if x != nil {
		return x.Spec
	}
	return nil
}

type GetTablesForSpec_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Marshalled []*schema.Table
	Tables []byte `protobuf:"bytes,3,opt,name=tables,proto3" json:"tables,omitempty"`
}

func (x *GetTablesForSpec_Response) Reset() {
	*x = GetTablesForSpec_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTablesForSpec_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTablesForSpec_Response) ProtoMessage() {}

func (x *GetTablesForSpec_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTablesForSpec_Response.ProtoReflect.Descriptor instead.
func (*GetTablesForSpec_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{4, 1}
}

func (x *GetTablesForSpec_Response) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetTablesForSpec_Response) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *GetTablesForSpec_Response) GetTables() []byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type GetSourceMetrics_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSourceMetrics_Request) Reset() {
	*x = GetSourceMetrics_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[16]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSourceMetrics_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSourceMetrics_Request) ProtoMessage() {}

func (x *GetSourceMetrics_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[16]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSourceMetrics_Request.ProtoReflect.Descriptor instead.
func (*GetSourceMetrics_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{5, 0}
}

type GetSourceMetrics_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled json of plugins.SourceMetrics
	Metrics []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *GetSourceMetrics_Response) Reset() {
	*x = GetSourceMetrics_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_source_proto_msgTypes[17]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSourceMetrics_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSourceMetrics_Response) ProtoMessage() {}

func (x *GetSourceMetrics_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_source_proto_msgTypes[17]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSourceMetrics_Response.ProtoReflect.Descriptor instead.
func (*GetSourceMetrics_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_source_proto_rawDescGZIP(), []int{5, 1}
}

func (x *GetSourceMetrics_Response) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_internal_pb_source_proto protoreflect.FileDescriptor

var file_internal_pb_source_proto_rawDesc = []byte{
	0x0a, 0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x62,
	0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x04, 0x53, 0x79, 0x6e,
	0x63, 0x1a, 0x1d, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x1a, 0x26, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x4e, 0x0a, 0x05, 0x53, 0x79, 0x6e, 0x63,
	0x32, 0x1a, 0x1d, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x1a, 0x26, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x41, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x53,
	0x79, 0x6e, 0x63, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x68, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x50, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x73, 0x46, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x1a, 0x1d, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x1a, 0x50, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x22, 0x43, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a,
	0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x32, 0x8e, 0x05, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x12, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x47, 0x65, 0x74,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x46, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x46, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x46, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x79, 0x6e, 0x63, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53,
	0x79, 0x6e, 0x63, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x36, 0x0a, 0x05, 0x53, 0x79, 0x6e, 0x63,
	0x32, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x32, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x79, 0x6e, 0x63, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01,
	0x12, 0x4f, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pb_source_proto_rawDescOnce sync.Once
	file_internal_pb_source_proto_rawDescData = file_internal_pb_source_proto_rawDesc
)

func file_internal_pb_source_proto_rawDescGZIP() []byte {
	file_internal_pb_source_proto_rawDescOnce.Do(func() {
		file_internal_pb_source_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pb_source_proto_rawDescData)
	})
	return file_internal_pb_source_proto_rawDescData
}

var file_internal_pb_source_proto_msgTypes = make([]protoimpl.MessageInfo, 18)
var file_internal_pb_source_proto_goTypes = []interface{}{
	(*Sync)(nil),                        // 0: proto.Sync
	(*Sync2)(nil),                       // 1: proto.Sync2
	(*GetSyncSummary)(nil),              // 2: proto.GetSyncSummary
	(*GetTables)(nil),                   // 3: proto.GetTables
	(*GetTablesForSpec)(nil),            // 4: proto.GetTablesForSpec
	(*GetSourceMetrics)(nil),            // 5: proto.GetSourceMetrics
	(*Sync_Request)(nil),                // 6: proto.Sync.Request
	(*Sync_Response)(nil),               // 7: proto.Sync.Response
	(*Sync2_Request)(nil),               // 8: proto.Sync2.Request
	(*Sync2_Response)(nil),              // 9: proto.Sync2.Response
	(*GetSyncSummary_Request)(nil),      // 10: proto.GetSyncSummary.Request
	(*GetSyncSummary_Response)(nil),     // 11: proto.GetSyncSummary.Response
	(*GetTables_Request)(nil),           // 12: proto.GetTables.Request
	(*GetTables_Response)(nil),          // 13: proto.GetTables.Response
	(*GetTablesForSpec_Request)(nil),    // 14: proto.GetTablesForSpec.Request
	(*GetTablesForSpec_Response)(nil),   // 15: proto.GetTablesForSpec.Response
	(*GetSourceMetrics_Request)(nil),    // 16: proto.GetSourceMetrics.Request
	(*GetSourceMetrics_Response)(nil),   // 17: proto.GetSourceMetrics.Response
	(*GetProtocolVersion_Request)(nil),  // 18: proto.GetProtocolVersion.Request
	(*GetName_Request)(nil),             // 19: proto.GetName.Request
	(*GetVersion_Request)(nil),          // 20: proto.GetVersion.Request
	(*GetProtocolVersion_Response)(nil), // 21: proto.GetProtocolVersion.Response
	(*GetName_Response)(nil),            // 22: proto.GetName.Response
	(*GetVersion_Response)(nil),         // 23: proto.GetVersion.Response
}
var file_internal_pb_source_proto_depIdxs = []int32{
	18, // 0: proto.Source.GetProtocolVersion:input_type -> proto.GetProtocolVersion.Request
	19, // 1: proto.Source.GetName:input_type -> proto.GetName.Request
	20, // 2: proto.Source.GetVersion:input_type -> proto.GetVersion.Request
	12, // 3: proto.Source.GetTables:input_type -> proto.GetTables.Request
	14, // 4: proto.Source.GetTablesForSpec:input_type -> proto.GetTablesForSpec.Request
	10, // 5: proto.Source.GetSyncSummary:input_type -> proto.GetSyncSummary.Request
	6,  // 6: proto.Source.Sync:input_type -> proto.Sync.Request
	8,  // 7: proto.Source.Sync2:input_type -> proto.Sync2.Request
	16, // 8: proto.Source.GetMetrics:input_type -> proto.GetSourceMetrics.Request
	21, // 9: proto.Source.GetProtocolVersion:output_type -> proto.GetProtocolVersion.Response
	22, // 10: proto.Source.GetName:output_type -> proto.GetName.Response
	23, // 11: proto.Source.GetVersion:output_type -> proto.GetVersion.Response
	13, // 12: proto.Source.GetTables:output_type -> proto.GetTables.Response
	15, // 13: proto.Source.GetTablesForSpec:output_type -> proto.GetTablesForSpec.Response
	11, // 14: proto.Source.GetSyncSummary:output_type -> proto.GetSyncSummary.Response
	7,  // 15: proto.Source.Sync:output_type -> proto.Sync.Response
	9,  // 16: proto.Source.Sync2:output_type -> proto.Sync2.Response
	17, // 17: proto.Source.GetMetrics:output_type -> proto.GetSourceMetrics.Response
	9,  // [9:18] is the sub-list for method output_type
	0,  // [0:9] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_internal_pb_source_proto_init() }
func file_internal_pb_source_proto_init() {
	if File_internal_pb_source_proto != nil {
		return
	}
	file_internal_pb_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_internal_pb_source_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync); i {
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
		file_internal_pb_source_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync2); i {
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
		file_internal_pb_source_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSyncSummary); i {
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
		file_internal_pb_source_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTables); i {
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
		file_internal_pb_source_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTablesForSpec); i {
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
		file_internal_pb_source_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSourceMetrics); i {
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
		file_internal_pb_source_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync_Request); i {
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
		file_internal_pb_source_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync_Response); i {
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
		file_internal_pb_source_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync2_Request); i {
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
		file_internal_pb_source_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sync2_Response); i {
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
		file_internal_pb_source_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSyncSummary_Request); i {
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
		file_internal_pb_source_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSyncSummary_Response); i {
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
		file_internal_pb_source_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTables_Request); i {
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
		file_internal_pb_source_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTables_Response); i {
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
		file_internal_pb_source_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTablesForSpec_Request); i {
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
		file_internal_pb_source_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTablesForSpec_Response); i {
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
		file_internal_pb_source_proto_msgTypes[16].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSourceMetrics_Request); i {
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
		file_internal_pb_source_proto_msgTypes[17].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSourceMetrics_Response); i {
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
			RawDescriptor: file_internal_pb_source_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   18,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pb_source_proto_goTypes,
		DependencyIndexes: file_internal_pb_source_proto_depIdxs,
		MessageInfos:      file_internal_pb_source_proto_msgTypes,
	}.Build()
	File_internal_pb_source_proto = out.File
	file_internal_pb_source_proto_rawDesc = nil
	file_internal_pb_source_proto_goTypes = nil
	file_internal_pb_source_proto_depIdxs = nil
}
