// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: pluginapi/v1/content.proto

package pluginapiv1

import (
	v1 "github.com/blackstork-io/fabric/plugin/ast/v1"
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

type LocationEffect int32

const (
	LocationEffect_LOCATION_EFFECT_UNSPECIFIED LocationEffect = 0
	LocationEffect_LOCATION_EFFECT_BEFORE      LocationEffect = 1
	LocationEffect_LOCATION_EFFECT_AFTER       LocationEffect = 2
)

// Enum value maps for LocationEffect.
var (
	LocationEffect_name = map[int32]string{
		0: "LOCATION_EFFECT_UNSPECIFIED",
		1: "LOCATION_EFFECT_BEFORE",
		2: "LOCATION_EFFECT_AFTER",
	}
	LocationEffect_value = map[string]int32{
		"LOCATION_EFFECT_UNSPECIFIED": 0,
		"LOCATION_EFFECT_BEFORE":      1,
		"LOCATION_EFFECT_AFTER":       2,
	}
)

func (x LocationEffect) Enum() *LocationEffect {
	p := new(LocationEffect)
	*p = x
	return p
}

func (x LocationEffect) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LocationEffect) Descriptor() protoreflect.EnumDescriptor {
	return file_pluginapi_v1_content_proto_enumTypes[0].Descriptor()
}

func (LocationEffect) Type() protoreflect.EnumType {
	return &file_pluginapi_v1_content_proto_enumTypes[0]
}

func (x LocationEffect) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LocationEffect.Descriptor instead.
func (LocationEffect) EnumDescriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{0}
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index  uint32         `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Effect LocationEffect `protobuf:"varint,2,opt,name=effect,proto3,enum=pluginapi.v1.LocationEffect" json:"effect,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{0}
}

func (x *Location) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Location) GetEffect() LocationEffect {
	if x != nil {
		return x.Effect
	}
	return LocationEffect_LOCATION_EFFECT_UNSPECIFIED
}

type ContentResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content  *Content  `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Location *Location `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *ContentResult) Reset() {
	*x = ContentResult{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentResult) ProtoMessage() {}

func (x *ContentResult) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentResult.ProtoReflect.Descriptor instead.
func (*ContentResult) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{1}
}

func (x *ContentResult) GetContent() *Content {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *ContentResult) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*Content_Element
	//	*Content_Section
	//	*Content_Empty
	Value isContent_Value `protobuf_oneof:"value"`
}

func (x *Content) Reset() {
	*x = Content{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{2}
}

func (m *Content) GetValue() isContent_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Content) GetElement() *ContentElement {
	if x, ok := x.GetValue().(*Content_Element); ok {
		return x.Element
	}
	return nil
}

func (x *Content) GetSection() *ContentSection {
	if x, ok := x.GetValue().(*Content_Section); ok {
		return x.Section
	}
	return nil
}

func (x *Content) GetEmpty() *ContentEmpty {
	if x, ok := x.GetValue().(*Content_Empty); ok {
		return x.Empty
	}
	return nil
}

type isContent_Value interface {
	isContent_Value()
}

type Content_Element struct {
	Element *ContentElement `protobuf:"bytes,1,opt,name=element,proto3,oneof"`
}

type Content_Section struct {
	Section *ContentSection `protobuf:"bytes,2,opt,name=section,proto3,oneof"`
}

type Content_Empty struct {
	Empty *ContentEmpty `protobuf:"bytes,3,opt,name=empty,proto3,oneof"`
}

func (*Content_Element) isContent_Value() {}

func (*Content_Section) isContent_Value() {}

func (*Content_Empty) isContent_Value() {}

type ContentSection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Children []*Content   `protobuf:"bytes,1,rep,name=children,proto3" json:"children,omitempty"`
	Meta     *v1.Metadata `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
}

func (x *ContentSection) Reset() {
	*x = ContentSection{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentSection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentSection) ProtoMessage() {}

func (x *ContentSection) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentSection.ProtoReflect.Descriptor instead.
func (*ContentSection) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{3}
}

func (x *ContentSection) GetChildren() []*Content {
	if x != nil {
		return x.Children
	}
	return nil
}

func (x *ContentSection) GetMeta() *v1.Metadata {
	if x != nil {
		return x.Meta
	}
	return nil
}

type ContentElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Markdown []byte                `protobuf:"bytes,1,opt,name=markdown,proto3" json:"markdown,omitempty"`
	Ast      *v1.FabricContentNode `protobuf:"bytes,2,opt,name=ast,proto3,oneof" json:"ast,omitempty"`
	Meta     *v1.Metadata          `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
}

func (x *ContentElement) Reset() {
	*x = ContentElement{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentElement) ProtoMessage() {}

func (x *ContentElement) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentElement.ProtoReflect.Descriptor instead.
func (*ContentElement) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{4}
}

func (x *ContentElement) GetMarkdown() []byte {
	if x != nil {
		return x.Markdown
	}
	return nil
}

func (x *ContentElement) GetAst() *v1.FabricContentNode {
	if x != nil {
		return x.Ast
	}
	return nil
}

func (x *ContentElement) GetMeta() *v1.Metadata {
	if x != nil {
		return x.Meta
	}
	return nil
}

type ContentEmpty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ContentEmpty) Reset() {
	*x = ContentEmpty{}
	mi := &file_pluginapi_v1_content_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentEmpty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentEmpty) ProtoMessage() {}

func (x *ContentEmpty) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_content_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentEmpty.ProtoReflect.Descriptor instead.
func (*ContentEmpty) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_content_proto_rawDescGZIP(), []int{5}
}

var File_pluginapi_v1_content_proto protoreflect.FileDescriptor

var file_pluginapi_v1_content_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x10, 0x61, 0x73, 0x74, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x08,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x34,
	0x0a, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x52, 0x06, 0x65, 0x66,
	0x66, 0x65, 0x63, 0x74, 0x22, 0x74, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2f, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xba, 0x01, 0x0a, 0x07, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x07, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x07, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x38, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x00, 0x52, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x07,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x69, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x08, 0x63, 0x68, 0x69,
	0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x04,
	0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x73, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65,
	0x74, 0x61, 0x22, 0x8c, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77,
	0x6e, 0x12, 0x30, 0x0a, 0x03, 0x61, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x61, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x62, 0x72, 0x69, 0x63, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x03, 0x61, 0x73, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x61, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x73,
	0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x2a, 0x68, 0x0a, 0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x12, 0x1f, 0x0a, 0x1b, 0x4c, 0x4f, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x45, 0x46, 0x46, 0x45, 0x43, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x4c, 0x4f, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x45, 0x46, 0x46, 0x45, 0x43, 0x54, 0x5f, 0x42, 0x45, 0x46, 0x4f, 0x52, 0x45, 0x10, 0x01,
	0x12, 0x19, 0x0a, 0x15, 0x4c, 0x4f, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x46, 0x46,
	0x45, 0x43, 0x54, 0x5f, 0x41, 0x46, 0x54, 0x45, 0x52, 0x10, 0x02, 0x42, 0xb2, 0x01, 0x0a, 0x10,
	0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x42, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x61,
	0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x6b, 0x2d, 0x69, 0x6f, 0x2f, 0x66, 0x61, 0x62, 0x72, 0x69,
	0x63, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61,
	0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pluginapi_v1_content_proto_rawDescOnce sync.Once
	file_pluginapi_v1_content_proto_rawDescData = file_pluginapi_v1_content_proto_rawDesc
)

func file_pluginapi_v1_content_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_content_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_content_proto_rawDescData = protoimpl.X.CompressGZIP(file_pluginapi_v1_content_proto_rawDescData)
	})
	return file_pluginapi_v1_content_proto_rawDescData
}

var file_pluginapi_v1_content_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pluginapi_v1_content_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pluginapi_v1_content_proto_goTypes = []any{
	(LocationEffect)(0),          // 0: pluginapi.v1.LocationEffect
	(*Location)(nil),             // 1: pluginapi.v1.Location
	(*ContentResult)(nil),        // 2: pluginapi.v1.ContentResult
	(*Content)(nil),              // 3: pluginapi.v1.Content
	(*ContentSection)(nil),       // 4: pluginapi.v1.ContentSection
	(*ContentElement)(nil),       // 5: pluginapi.v1.ContentElement
	(*ContentEmpty)(nil),         // 6: pluginapi.v1.ContentEmpty
	(*v1.Metadata)(nil),          // 7: ast.v1.Metadata
	(*v1.FabricContentNode)(nil), // 8: ast.v1.FabricContentNode
}
var file_pluginapi_v1_content_proto_depIdxs = []int32{
	0,  // 0: pluginapi.v1.Location.effect:type_name -> pluginapi.v1.LocationEffect
	3,  // 1: pluginapi.v1.ContentResult.content:type_name -> pluginapi.v1.Content
	1,  // 2: pluginapi.v1.ContentResult.location:type_name -> pluginapi.v1.Location
	5,  // 3: pluginapi.v1.Content.element:type_name -> pluginapi.v1.ContentElement
	4,  // 4: pluginapi.v1.Content.section:type_name -> pluginapi.v1.ContentSection
	6,  // 5: pluginapi.v1.Content.empty:type_name -> pluginapi.v1.ContentEmpty
	3,  // 6: pluginapi.v1.ContentSection.children:type_name -> pluginapi.v1.Content
	7,  // 7: pluginapi.v1.ContentSection.meta:type_name -> ast.v1.Metadata
	8,  // 8: pluginapi.v1.ContentElement.ast:type_name -> ast.v1.FabricContentNode
	7,  // 9: pluginapi.v1.ContentElement.meta:type_name -> ast.v1.Metadata
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_content_proto_init() }
func file_pluginapi_v1_content_proto_init() {
	if File_pluginapi_v1_content_proto != nil {
		return
	}
	file_pluginapi_v1_content_proto_msgTypes[2].OneofWrappers = []any{
		(*Content_Element)(nil),
		(*Content_Section)(nil),
		(*Content_Empty)(nil),
	}
	file_pluginapi_v1_content_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pluginapi_v1_content_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_content_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_content_proto_depIdxs,
		EnumInfos:         file_pluginapi_v1_content_proto_enumTypes,
		MessageInfos:      file_pluginapi_v1_content_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_content_proto = out.File
	file_pluginapi_v1_content_proto_rawDesc = nil
	file_pluginapi_v1_content_proto_goTypes = nil
	file_pluginapi_v1_content_proto_depIdxs = nil
}
