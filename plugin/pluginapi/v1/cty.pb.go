// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        (unknown)
// source: pluginapi/v1/cty.proto

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

type Cty struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Cty with nil data is decoded as cty.NilVal
	//
	// Types that are valid to be assigned to Data:
	//
	//	*Cty_Primitive_
	//	*Cty_Object_
	//	*Cty_Map
	//	*Cty_List
	//	*Cty_Set
	//	*Cty_Tuple
	//	*Cty_Null
	//	*Cty_Caps
	//	*Cty_Unknown
	//	*Cty_Dyn
	Data          isCty_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty) Reset() {
	*x = Cty{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty) ProtoMessage() {}

func (x *Cty) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty.ProtoReflect.Descriptor instead.
func (*Cty) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0}
}

func (x *Cty) GetData() isCty_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty) GetPrimitive() *Cty_Primitive {
	if x != nil {
		if x, ok := x.Data.(*Cty_Primitive_); ok {
			return x.Primitive
		}
	}
	return nil
}

func (x *Cty) GetObject() *Cty_Object {
	if x != nil {
		if x, ok := x.Data.(*Cty_Object_); ok {
			return x.Object
		}
	}
	return nil
}

func (x *Cty) GetMap() *Cty_Mapping {
	if x != nil {
		if x, ok := x.Data.(*Cty_Map); ok {
			return x.Map
		}
	}
	return nil
}

func (x *Cty) GetList() *Cty_Sequence {
	if x != nil {
		if x, ok := x.Data.(*Cty_List); ok {
			return x.List
		}
	}
	return nil
}

func (x *Cty) GetSet() *Cty_Sequence {
	if x != nil {
		if x, ok := x.Data.(*Cty_Set); ok {
			return x.Set
		}
	}
	return nil
}

func (x *Cty) GetTuple() *Cty_Sequence {
	if x != nil {
		if x, ok := x.Data.(*Cty_Tuple); ok {
			return x.Tuple
		}
	}
	return nil
}

func (x *Cty) GetNull() *CtyType {
	if x != nil {
		if x, ok := x.Data.(*Cty_Null); ok {
			return x.Null
		}
	}
	return nil
}

func (x *Cty) GetCaps() *Cty_Capsule {
	if x != nil {
		if x, ok := x.Data.(*Cty_Caps); ok {
			return x.Caps
		}
	}
	return nil
}

func (x *Cty) GetUnknown() *CtyType {
	if x != nil {
		if x, ok := x.Data.(*Cty_Unknown); ok {
			return x.Unknown
		}
	}
	return nil
}

func (x *Cty) GetDyn() *Cty_Dynamic {
	if x != nil {
		if x, ok := x.Data.(*Cty_Dyn); ok {
			return x.Dyn
		}
	}
	return nil
}

type isCty_Data interface {
	isCty_Data()
}

type Cty_Primitive_ struct {
	Primitive *Cty_Primitive `protobuf:"bytes,1,opt,name=primitive,proto3,oneof"`
}

type Cty_Object_ struct {
	Object *Cty_Object `protobuf:"bytes,2,opt,name=object,proto3,oneof"`
}

type Cty_Map struct {
	Map *Cty_Mapping `protobuf:"bytes,3,opt,name=map,proto3,oneof"`
}

type Cty_List struct {
	List *Cty_Sequence `protobuf:"bytes,4,opt,name=list,proto3,oneof"`
}

type Cty_Set struct {
	Set *Cty_Sequence `protobuf:"bytes,5,opt,name=set,proto3,oneof"`
}

type Cty_Tuple struct {
	Tuple *Cty_Sequence `protobuf:"bytes,6,opt,name=tuple,proto3,oneof"`
}

type Cty_Null struct {
	// Specifies type of null value
	Null *CtyType `protobuf:"bytes,7,opt,name=null,proto3,oneof"`
}

type Cty_Caps struct {
	Caps *Cty_Capsule `protobuf:"bytes,8,opt,name=caps,proto3,oneof"`
}

type Cty_Unknown struct {
	// Specifies type of the unknown value
	Unknown *CtyType `protobuf:"bytes,9,opt,name=unknown,proto3,oneof"`
}

type Cty_Dyn struct {
	// DynamicPseudoType
	Dyn *Cty_Dynamic `protobuf:"bytes,10,opt,name=dyn,proto3,oneof"`
}

func (*Cty_Primitive_) isCty_Data() {}

func (*Cty_Object_) isCty_Data() {}

func (*Cty_Map) isCty_Data() {}

func (*Cty_List) isCty_Data() {}

func (*Cty_Set) isCty_Data() {}

func (*Cty_Tuple) isCty_Data() {}

func (*Cty_Null) isCty_Data() {}

func (*Cty_Caps) isCty_Data() {}

func (*Cty_Unknown) isCty_Data() {}

func (*Cty_Dyn) isCty_Data() {}

// Forces decoding of the inner Cty as a type
type CtyType struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          *Cty                   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CtyType) Reset() {
	*x = CtyType{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CtyType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyType) ProtoMessage() {}

func (x *CtyType) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyType.ProtoReflect.Descriptor instead.
func (*CtyType) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{1}
}

func (x *CtyType) GetType() *Cty {
	if x != nil {
		return x.Type
	}
	return nil
}

// Forces decoding of the inner Cty as a value
type CtyValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         *Cty                   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CtyValue) Reset() {
	*x = CtyValue{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CtyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyValue) ProtoMessage() {}

func (x *CtyValue) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyValue.ProtoReflect.Descriptor instead.
func (*CtyValue) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{2}
}

func (x *CtyValue) GetValue() *Cty {
	if x != nil {
		return x.Value
	}
	return nil
}

type Cty_Primitive struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*Cty_Primitive_Str
	//	*Cty_Primitive_Num
	//	*Cty_Primitive_Bln
	Data          isCty_Primitive_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Primitive) Reset() {
	*x = Cty_Primitive{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Primitive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Primitive) ProtoMessage() {}

func (x *Cty_Primitive) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Primitive.ProtoReflect.Descriptor instead.
func (*Cty_Primitive) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Cty_Primitive) GetData() isCty_Primitive_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty_Primitive) GetStr() string {
	if x != nil {
		if x, ok := x.Data.(*Cty_Primitive_Str); ok {
			return x.Str
		}
	}
	return ""
}

func (x *Cty_Primitive) GetNum() []byte {
	if x != nil {
		if x, ok := x.Data.(*Cty_Primitive_Num); ok {
			return x.Num
		}
	}
	return nil
}

func (x *Cty_Primitive) GetBln() bool {
	if x != nil {
		if x, ok := x.Data.(*Cty_Primitive_Bln); ok {
			return x.Bln
		}
	}
	return false
}

type isCty_Primitive_Data interface {
	isCty_Primitive_Data()
}

type Cty_Primitive_Str struct {
	Str string `protobuf:"bytes,1,opt,name=str,proto3,oneof"`
}

type Cty_Primitive_Num struct {
	// empty value is used for marking the type
	Num []byte `protobuf:"bytes,2,opt,name=num,proto3,oneof"`
}

type Cty_Primitive_Bln struct {
	Bln bool `protobuf:"varint,3,opt,name=bln,proto3,oneof"`
}

func (*Cty_Primitive_Str) isCty_Primitive_Data() {}

func (*Cty_Primitive_Num) isCty_Primitive_Data() {}

func (*Cty_Primitive_Bln) isCty_Primitive_Data() {}

type Cty_Object struct {
	state         protoimpl.MessageState      `protogen:"open.v1"`
	Data          map[string]*Cty_Object_Attr `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Object) Reset() {
	*x = Cty_Object{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Object) ProtoMessage() {}

func (x *Cty_Object) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Object.ProtoReflect.Descriptor instead.
func (*Cty_Object) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Cty_Object) GetData() map[string]*Cty_Object_Attr {
	if x != nil {
		return x.Data
	}
	return nil
}

type Cty_Mapping struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Data  map[string]*Cty        `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Original map is empty, element was added to preserve the type
	OnlyType      bool `protobuf:"varint,2,opt,name=only_type,json=onlyType,proto3" json:"only_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Mapping) Reset() {
	*x = Cty_Mapping{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Mapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Mapping) ProtoMessage() {}

func (x *Cty_Mapping) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Mapping.ProtoReflect.Descriptor instead.
func (*Cty_Mapping) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Cty_Mapping) GetData() map[string]*Cty {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty_Mapping) GetOnlyType() bool {
	if x != nil {
		return x.OnlyType
	}
	return false
}

type Cty_Sequence struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Data  []*Cty                 `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	// Original sequence is empty, element added to preserve the type
	// Not true for empty tuples, since they are valid values
	OnlyType      bool `protobuf:"varint,2,opt,name=only_type,json=onlyType,proto3" json:"only_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Sequence) Reset() {
	*x = Cty_Sequence{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Sequence) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Sequence) ProtoMessage() {}

func (x *Cty_Sequence) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Sequence.ProtoReflect.Descriptor instead.
func (*Cty_Sequence) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 3}
}

func (x *Cty_Sequence) GetData() []*Cty {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty_Sequence) GetOnlyType() bool {
	if x != nil {
		return x.OnlyType
	}
	return false
}

type Cty_Capsule struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*Cty_Capsule_PluginData
	Data          isCty_Capsule_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Capsule) Reset() {
	*x = Cty_Capsule{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Capsule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Capsule) ProtoMessage() {}

func (x *Cty_Capsule) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Capsule.ProtoReflect.Descriptor instead.
func (*Cty_Capsule) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 4}
}

func (x *Cty_Capsule) GetData() isCty_Capsule_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty_Capsule) GetPluginData() *Data {
	if x != nil {
		if x, ok := x.Data.(*Cty_Capsule_PluginData); ok {
			return x.PluginData
		}
	}
	return nil
}

type isCty_Capsule_Data interface {
	isCty_Capsule_Data()
}

type Cty_Capsule_PluginData struct {
	PluginData *Data `protobuf:"bytes,1,opt,name=plugin_data,json=pluginData,proto3,oneof"`
}

func (*Cty_Capsule_PluginData) isCty_Capsule_Data() {}

type Cty_Dynamic struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Dynamic) Reset() {
	*x = Cty_Dynamic{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Dynamic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Dynamic) ProtoMessage() {}

func (x *Cty_Dynamic) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Dynamic.ProtoReflect.Descriptor instead.
func (*Cty_Dynamic) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 5}
}

type Cty_Object_Attr struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *Cty                   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Optional      bool                   `protobuf:"varint,2,opt,name=optional,proto3" json:"optional,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cty_Object_Attr) Reset() {
	*x = Cty_Object_Attr{}
	mi := &file_pluginapi_v1_cty_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cty_Object_Attr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cty_Object_Attr) ProtoMessage() {}

func (x *Cty_Object_Attr) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_cty_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cty_Object_Attr.ProtoReflect.Descriptor instead.
func (*Cty_Object_Attr) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_cty_proto_rawDescGZIP(), []int{0, 1, 0}
}

func (x *Cty_Object_Attr) GetData() *Cty {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Cty_Object_Attr) GetOptional() bool {
	if x != nil {
		return x.Optional
	}
	return false
}

var File_pluginapi_v1_cty_proto protoreflect.FileDescriptor

var file_pluginapi_v1_cty_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8d, 0x09, 0x0a, 0x03, 0x43, 0x74, 0x79, 0x12, 0x3b, 0x0a, 0x09, 0x70, 0x72, 0x69, 0x6d, 0x69,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x50, 0x72,
	0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x48, 0x00, 0x52, 0x09, 0x70, 0x72, 0x69, 0x6d, 0x69,
	0x74, 0x69, 0x76, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x48, 0x00,
	0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x2d, 0x0a, 0x03, 0x6d, 0x61, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x48, 0x00, 0x52, 0x03, 0x6d, 0x61, 0x70, 0x12, 0x30, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63,
	0x65, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x03, 0x73, 0x65, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x48, 0x00, 0x52, 0x03, 0x73, 0x65, 0x74, 0x12, 0x32, 0x0a, 0x05, 0x74, 0x75, 0x70,
	0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x48, 0x00, 0x52, 0x05, 0x74, 0x75, 0x70, 0x6c, 0x65, 0x12, 0x2b, 0x0a,
	0x04, 0x6e, 0x75, 0x6c, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x54, 0x79,
	0x70, 0x65, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x75, 0x6c, 0x6c, 0x12, 0x2f, 0x0a, 0x04, 0x63, 0x61,
	0x70, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x43, 0x61, 0x70, 0x73,
	0x75, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x04, 0x63, 0x61, 0x70, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x75,
	0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x48, 0x00, 0x52, 0x07, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x12, 0x2d,
	0x0a, 0x03, 0x64, 0x79, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x44,
	0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x48, 0x00, 0x52, 0x03, 0x64, 0x79, 0x6e, 0x1a, 0x4f, 0x0a,
	0x09, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x12, 0x12, 0x0a, 0x03, 0x73, 0x74,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x73, 0x74, 0x72, 0x12, 0x12,
	0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x03, 0x6e,
	0x75, 0x6d, 0x12, 0x12, 0x0a, 0x03, 0x62, 0x6c, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x00, 0x52, 0x03, 0x62, 0x6c, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0xe3,
	0x01, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x1a, 0x49, 0x0a, 0x04, 0x41, 0x74, 0x74, 0x72, 0x12, 0x25, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x1a, 0x56, 0x0a, 0x09,
	0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x2e, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0xab, 0x01, 0x0a, 0x07, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x12, 0x37, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74,
	0x79, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x6e, 0x6c,
	0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x6e,
	0x6c, 0x79, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x4a, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x1a, 0x4e, 0x0a, 0x08, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x25,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x6e, 0x6c, 0x79, 0x54, 0x79,
	0x70, 0x65, 0x1a, 0x48, 0x0a, 0x07, 0x43, 0x61, 0x70, 0x73, 0x75, 0x6c, 0x65, 0x12, 0x35, 0x0a,
	0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x44, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x09, 0x0a, 0x07,
	0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x30, 0x0a, 0x07, 0x43, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x22, 0x33, 0x0a, 0x08, 0x43, 0x74, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x27, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0xae, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x08, 0x43, 0x74, 0x79,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x6b, 0x2d, 0x69,
	0x6f, 0x2f, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02,
	0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c,
	0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pluginapi_v1_cty_proto_rawDescOnce sync.Once
	file_pluginapi_v1_cty_proto_rawDescData = file_pluginapi_v1_cty_proto_rawDesc
)

func file_pluginapi_v1_cty_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_cty_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_cty_proto_rawDescData = protoimpl.X.CompressGZIP(file_pluginapi_v1_cty_proto_rawDescData)
	})
	return file_pluginapi_v1_cty_proto_rawDescData
}

var file_pluginapi_v1_cty_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_pluginapi_v1_cty_proto_goTypes = []any{
	(*Cty)(nil),             // 0: pluginapi.v1.Cty
	(*CtyType)(nil),         // 1: pluginapi.v1.CtyType
	(*CtyValue)(nil),        // 2: pluginapi.v1.CtyValue
	(*Cty_Primitive)(nil),   // 3: pluginapi.v1.Cty.Primitive
	(*Cty_Object)(nil),      // 4: pluginapi.v1.Cty.Object
	(*Cty_Mapping)(nil),     // 5: pluginapi.v1.Cty.Mapping
	(*Cty_Sequence)(nil),    // 6: pluginapi.v1.Cty.Sequence
	(*Cty_Capsule)(nil),     // 7: pluginapi.v1.Cty.Capsule
	(*Cty_Dynamic)(nil),     // 8: pluginapi.v1.Cty.Dynamic
	(*Cty_Object_Attr)(nil), // 9: pluginapi.v1.Cty.Object.Attr
	nil,                     // 10: pluginapi.v1.Cty.Object.DataEntry
	nil,                     // 11: pluginapi.v1.Cty.Mapping.DataEntry
	(*Data)(nil),            // 12: pluginapi.v1.Data
}
var file_pluginapi_v1_cty_proto_depIdxs = []int32{
	3,  // 0: pluginapi.v1.Cty.primitive:type_name -> pluginapi.v1.Cty.Primitive
	4,  // 1: pluginapi.v1.Cty.object:type_name -> pluginapi.v1.Cty.Object
	5,  // 2: pluginapi.v1.Cty.map:type_name -> pluginapi.v1.Cty.Mapping
	6,  // 3: pluginapi.v1.Cty.list:type_name -> pluginapi.v1.Cty.Sequence
	6,  // 4: pluginapi.v1.Cty.set:type_name -> pluginapi.v1.Cty.Sequence
	6,  // 5: pluginapi.v1.Cty.tuple:type_name -> pluginapi.v1.Cty.Sequence
	1,  // 6: pluginapi.v1.Cty.null:type_name -> pluginapi.v1.CtyType
	7,  // 7: pluginapi.v1.Cty.caps:type_name -> pluginapi.v1.Cty.Capsule
	1,  // 8: pluginapi.v1.Cty.unknown:type_name -> pluginapi.v1.CtyType
	8,  // 9: pluginapi.v1.Cty.dyn:type_name -> pluginapi.v1.Cty.Dynamic
	0,  // 10: pluginapi.v1.CtyType.type:type_name -> pluginapi.v1.Cty
	0,  // 11: pluginapi.v1.CtyValue.value:type_name -> pluginapi.v1.Cty
	10, // 12: pluginapi.v1.Cty.Object.data:type_name -> pluginapi.v1.Cty.Object.DataEntry
	11, // 13: pluginapi.v1.Cty.Mapping.data:type_name -> pluginapi.v1.Cty.Mapping.DataEntry
	0,  // 14: pluginapi.v1.Cty.Sequence.data:type_name -> pluginapi.v1.Cty
	12, // 15: pluginapi.v1.Cty.Capsule.plugin_data:type_name -> pluginapi.v1.Data
	0,  // 16: pluginapi.v1.Cty.Object.Attr.data:type_name -> pluginapi.v1.Cty
	9,  // 17: pluginapi.v1.Cty.Object.DataEntry.value:type_name -> pluginapi.v1.Cty.Object.Attr
	0,  // 18: pluginapi.v1.Cty.Mapping.DataEntry.value:type_name -> pluginapi.v1.Cty
	19, // [19:19] is the sub-list for method output_type
	19, // [19:19] is the sub-list for method input_type
	19, // [19:19] is the sub-list for extension type_name
	19, // [19:19] is the sub-list for extension extendee
	0,  // [0:19] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_cty_proto_init() }
func file_pluginapi_v1_cty_proto_init() {
	if File_pluginapi_v1_cty_proto != nil {
		return
	}
	file_pluginapi_v1_data_proto_init()
	file_pluginapi_v1_cty_proto_msgTypes[0].OneofWrappers = []any{
		(*Cty_Primitive_)(nil),
		(*Cty_Object_)(nil),
		(*Cty_Map)(nil),
		(*Cty_List)(nil),
		(*Cty_Set)(nil),
		(*Cty_Tuple)(nil),
		(*Cty_Null)(nil),
		(*Cty_Caps)(nil),
		(*Cty_Unknown)(nil),
		(*Cty_Dyn)(nil),
	}
	file_pluginapi_v1_cty_proto_msgTypes[3].OneofWrappers = []any{
		(*Cty_Primitive_Str)(nil),
		(*Cty_Primitive_Num)(nil),
		(*Cty_Primitive_Bln)(nil),
	}
	file_pluginapi_v1_cty_proto_msgTypes[7].OneofWrappers = []any{
		(*Cty_Capsule_PluginData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pluginapi_v1_cty_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_cty_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_cty_proto_depIdxs,
		MessageInfos:      file_pluginapi_v1_cty_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_cty_proto = out.File
	file_pluginapi_v1_cty_proto_rawDesc = nil
	file_pluginapi_v1_cty_proto_goTypes = nil
	file_pluginapi_v1_cty_proto_depIdxs = nil
}
