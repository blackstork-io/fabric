// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: pluginapi/v1/data.proto

package pluginapiv1

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

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*Data_StringVal
	//	*Data_NumberVal
	//	*Data_BoolVal
	//	*Data_ListVal
	//	*Data_MapVal
	Data isData_Data `protobuf_oneof:"data"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{0}
}

func (m *Data) GetData() isData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Data) GetStringVal() *StringData {
	if x, ok := x.GetData().(*Data_StringVal); ok {
		return x.StringVal
	}
	return nil
}

func (x *Data) GetNumberVal() *NumberData {
	if x, ok := x.GetData().(*Data_NumberVal); ok {
		return x.NumberVal
	}
	return nil
}

func (x *Data) GetBoolVal() *BoolData {
	if x, ok := x.GetData().(*Data_BoolVal); ok {
		return x.BoolVal
	}
	return nil
}

func (x *Data) GetListVal() *ListData {
	if x, ok := x.GetData().(*Data_ListVal); ok {
		return x.ListVal
	}
	return nil
}

func (x *Data) GetMapVal() *MapData {
	if x, ok := x.GetData().(*Data_MapVal); ok {
		return x.MapVal
	}
	return nil
}

type isData_Data interface {
	isData_Data()
}

type Data_StringVal struct {
	StringVal *StringData `protobuf:"bytes,1,opt,name=string_val,json=stringVal,proto3,oneof"`
}

type Data_NumberVal struct {
	NumberVal *NumberData `protobuf:"bytes,2,opt,name=number_val,json=numberVal,proto3,oneof"`
}

type Data_BoolVal struct {
	BoolVal *BoolData `protobuf:"bytes,3,opt,name=bool_val,json=boolVal,proto3,oneof"`
}

type Data_ListVal struct {
	ListVal *ListData `protobuf:"bytes,4,opt,name=list_val,json=listVal,proto3,oneof"`
}

type Data_MapVal struct {
	MapVal *MapData `protobuf:"bytes,5,opt,name=map_val,json=mapVal,proto3,oneof"`
}

func (*Data_StringVal) isData_Data() {}

func (*Data_NumberVal) isData_Data() {}

func (*Data_BoolVal) isData_Data() {}

func (*Data_ListVal) isData_Data() {}

func (*Data_MapVal) isData_Data() {}

type ListData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*Data `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *ListData) Reset() {
	*x = ListData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListData) ProtoMessage() {}

func (x *ListData) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListData.ProtoReflect.Descriptor instead.
func (*ListData) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{1}
}

func (x *ListData) GetValue() []*Data {
	if x != nil {
		return x.Value
	}
	return nil
}

type MapData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value map[string]*Data `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MapData) Reset() {
	*x = MapData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapData) ProtoMessage() {}

func (x *MapData) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapData.ProtoReflect.Descriptor instead.
func (*MapData) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{2}
}

func (x *MapData) GetValue() map[string]*Data {
	if x != nil {
		return x.Value
	}
	return nil
}

type StringData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringData) Reset() {
	*x = StringData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringData) ProtoMessage() {}

func (x *StringData) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringData.ProtoReflect.Descriptor instead.
func (*StringData) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{3}
}

func (x *StringData) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type NumberData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *NumberData) Reset() {
	*x = NumberData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NumberData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NumberData) ProtoMessage() {}

func (x *NumberData) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NumberData.ProtoReflect.Descriptor instead.
func (*NumberData) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{4}
}

func (x *NumberData) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type BoolData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value bool `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BoolData) Reset() {
	*x = BoolData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_v1_data_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolData) ProtoMessage() {}

func (x *BoolData) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_data_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolData.ProtoReflect.Descriptor instead.
func (*BoolData) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_data_proto_rawDescGZIP(), []int{5}
}

func (x *BoolData) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

var File_pluginapi_v1_data_proto protoreflect.FileDescriptor

var file_pluginapi_v1_data_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x22, 0xa0, 0x02, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00,
	0x52, 0x09, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x12, 0x33, 0x0a, 0x08, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x76,
	0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x44, 0x61, 0x74, 0x61,
	0x48, 0x00, 0x52, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x12, 0x33, 0x0a, 0x08, 0x6c,
	0x69, 0x73, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x07, 0x6c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c,
	0x12, 0x30, 0x0a, 0x07, 0x6d, 0x61, 0x70, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x4d, 0x61, 0x70, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x06, 0x6d, 0x61, 0x70, 0x56,
	0x61, 0x6c, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x34, 0x0a, 0x08, 0x4c, 0x69,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x8f, 0x01, 0x0a, 0x07, 0x4d, 0x61, 0x70, 0x44, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x70, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x1a, 0x4c, 0x0a, 0x0a, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x22, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a, 0x08, 0x42, 0x6f,
	0x6f, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0xaf, 0x01, 0x0a,
	0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x42, 0x09, 0x44, 0x61, 0x74, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b,
	0x73, 0x74, 0x6f, 0x72, 0x6b, 0x2d, 0x69, 0x6f, 0x2f, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x2f,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x3b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pluginapi_v1_data_proto_rawDescOnce sync.Once
	file_pluginapi_v1_data_proto_rawDescData = file_pluginapi_v1_data_proto_rawDesc
)

func file_pluginapi_v1_data_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_data_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_pluginapi_v1_data_proto_rawDescData)
	})
	return file_pluginapi_v1_data_proto_rawDescData
}

var file_pluginapi_v1_data_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pluginapi_v1_data_proto_goTypes = []interface{}{
	(*Data)(nil),       // 0: pluginapi.v1.Data
	(*ListData)(nil),   // 1: pluginapi.v1.ListData
	(*MapData)(nil),    // 2: pluginapi.v1.MapData
	(*StringData)(nil), // 3: pluginapi.v1.StringData
	(*NumberData)(nil), // 4: pluginapi.v1.NumberData
	(*BoolData)(nil),   // 5: pluginapi.v1.BoolData
	nil,                // 6: pluginapi.v1.MapData.ValueEntry
}
var file_pluginapi_v1_data_proto_depIdxs = []int32{
	3, // 0: pluginapi.v1.Data.string_val:type_name -> pluginapi.v1.StringData
	4, // 1: pluginapi.v1.Data.number_val:type_name -> pluginapi.v1.NumberData
	5, // 2: pluginapi.v1.Data.bool_val:type_name -> pluginapi.v1.BoolData
	1, // 3: pluginapi.v1.Data.list_val:type_name -> pluginapi.v1.ListData
	2, // 4: pluginapi.v1.Data.map_val:type_name -> pluginapi.v1.MapData
	0, // 5: pluginapi.v1.ListData.value:type_name -> pluginapi.v1.Data
	6, // 6: pluginapi.v1.MapData.value:type_name -> pluginapi.v1.MapData.ValueEntry
	0, // 7: pluginapi.v1.MapData.ValueEntry.value:type_name -> pluginapi.v1.Data
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_data_proto_init() }
func file_pluginapi_v1_data_proto_init() {
	if File_pluginapi_v1_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pluginapi_v1_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_pluginapi_v1_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListData); i {
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
		file_pluginapi_v1_data_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapData); i {
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
		file_pluginapi_v1_data_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringData); i {
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
		file_pluginapi_v1_data_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NumberData); i {
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
		file_pluginapi_v1_data_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolData); i {
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
	file_pluginapi_v1_data_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Data_StringVal)(nil),
		(*Data_NumberVal)(nil),
		(*Data_BoolVal)(nil),
		(*Data_ListVal)(nil),
		(*Data_MapVal)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pluginapi_v1_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_data_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_data_proto_depIdxs,
		MessageInfos:      file_pluginapi_v1_data_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_data_proto = out.File
	file_pluginapi_v1_data_proto_rawDesc = nil
	file_pluginapi_v1_data_proto_goTypes = nil
	file_pluginapi_v1_data_proto_depIdxs = nil
}