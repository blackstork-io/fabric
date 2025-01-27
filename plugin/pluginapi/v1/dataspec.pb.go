// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: pluginapi/v1/dataspec.proto

package pluginapiv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AttrSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type          *CtyType               `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	DefaultVal    *CtyValue              `protobuf:"bytes,3,opt,name=default_val,json=defaultVal,proto3" json:"default_val,omitempty"`
	ExampleVal    *CtyValue              `protobuf:"bytes,4,opt,name=example_val,json=exampleVal,proto3" json:"example_val,omitempty"`
	Doc           string                 `protobuf:"bytes,5,opt,name=doc,proto3" json:"doc,omitempty"`
	Constraints   uint32                 `protobuf:"varint,7,opt,name=constraints,proto3" json:"constraints,omitempty"`
	OneOf         []*CtyValue            `protobuf:"bytes,8,rep,name=one_of,json=oneOf,proto3" json:"one_of,omitempty"`
	MinInclusive  *CtyValue              `protobuf:"bytes,9,opt,name=min_inclusive,json=minInclusive,proto3" json:"min_inclusive,omitempty"`
	MaxInclusive  *CtyValue              `protobuf:"bytes,10,opt,name=max_inclusive,json=maxInclusive,proto3" json:"max_inclusive,omitempty"`
	Deprecated    string                 `protobuf:"bytes,11,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	Secret        bool                   `protobuf:"varint,12,opt,name=secret,proto3" json:"secret,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AttrSpec) Reset() {
	*x = AttrSpec{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AttrSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttrSpec) ProtoMessage() {}

func (x *AttrSpec) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttrSpec.ProtoReflect.Descriptor instead.
func (*AttrSpec) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{0}
}

func (x *AttrSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AttrSpec) GetType() *CtyType {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *AttrSpec) GetDefaultVal() *CtyValue {
	if x != nil {
		return x.DefaultVal
	}
	return nil
}

func (x *AttrSpec) GetExampleVal() *CtyValue {
	if x != nil {
		return x.ExampleVal
	}
	return nil
}

func (x *AttrSpec) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *AttrSpec) GetConstraints() uint32 {
	if x != nil {
		return x.Constraints
	}
	return 0
}

func (x *AttrSpec) GetOneOf() []*CtyValue {
	if x != nil {
		return x.OneOf
	}
	return nil
}

func (x *AttrSpec) GetMinInclusive() *CtyValue {
	if x != nil {
		return x.MinInclusive
	}
	return nil
}

func (x *AttrSpec) GetMaxInclusive() *CtyValue {
	if x != nil {
		return x.MaxInclusive
	}
	return nil
}

func (x *AttrSpec) GetDeprecated() string {
	if x != nil {
		return x.Deprecated
	}
	return ""
}

func (x *AttrSpec) GetSecret() bool {
	if x != nil {
		return x.Secret
	}
	return false
}

type BlockSpec struct {
	state                      protoimpl.MessageState   `protogen:"open.v1"`
	HeadersSpec                []*BlockSpec_NameMatcher `protobuf:"bytes,1,rep,name=headers_spec,json=headersSpec,proto3" json:"headers_spec,omitempty"`
	Required                   bool                     `protobuf:"varint,2,opt,name=required,proto3" json:"required,omitempty"`
	Repeatable                 bool                     `protobuf:"varint,3,opt,name=repeatable,proto3" json:"repeatable,omitempty"`
	Doc                        string                   `protobuf:"bytes,4,opt,name=doc,proto3" json:"doc,omitempty"`
	BlockSpecs                 []*BlockSpec             `protobuf:"bytes,5,rep,name=block_specs,json=blockSpecs,proto3" json:"block_specs,omitempty"`
	AttrSpecs                  []*AttrSpec              `protobuf:"bytes,6,rep,name=attr_specs,json=attrSpecs,proto3" json:"attr_specs,omitempty"`
	AllowUnspecifiedBlocks     bool                     `protobuf:"varint,7,opt,name=allow_unspecified_blocks,json=allowUnspecifiedBlocks,proto3" json:"allow_unspecified_blocks,omitempty"`
	AllowUnspecifiedAttributes bool                     `protobuf:"varint,8,opt,name=allow_unspecified_attributes,json=allowUnspecifiedAttributes,proto3" json:"allow_unspecified_attributes,omitempty"`
	unknownFields              protoimpl.UnknownFields
	sizeCache                  protoimpl.SizeCache
}

func (x *BlockSpec) Reset() {
	*x = BlockSpec{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockSpec) ProtoMessage() {}

func (x *BlockSpec) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockSpec.ProtoReflect.Descriptor instead.
func (*BlockSpec) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{1}
}

func (x *BlockSpec) GetHeadersSpec() []*BlockSpec_NameMatcher {
	if x != nil {
		return x.HeadersSpec
	}
	return nil
}

func (x *BlockSpec) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

func (x *BlockSpec) GetRepeatable() bool {
	if x != nil {
		return x.Repeatable
	}
	return false
}

func (x *BlockSpec) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *BlockSpec) GetBlockSpecs() []*BlockSpec {
	if x != nil {
		return x.BlockSpecs
	}
	return nil
}

func (x *BlockSpec) GetAttrSpecs() []*AttrSpec {
	if x != nil {
		return x.AttrSpecs
	}
	return nil
}

func (x *BlockSpec) GetAllowUnspecifiedBlocks() bool {
	if x != nil {
		return x.AllowUnspecifiedBlocks
	}
	return false
}

func (x *BlockSpec) GetAllowUnspecifiedAttributes() bool {
	if x != nil {
		return x.AllowUnspecifiedAttributes
	}
	return false
}

type Attr struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	NameRange     *Range                 `protobuf:"bytes,2,opt,name=name_range,json=nameRange,proto3" json:"name_range,omitempty"`
	Value         *CtyValue              `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	ValueRange    *Range                 `protobuf:"bytes,4,opt,name=value_range,json=valueRange,proto3" json:"value_range,omitempty"`
	Secret        bool                   `protobuf:"varint,5,opt,name=secret,proto3" json:"secret,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Attr) Reset() {
	*x = Attr{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Attr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attr) ProtoMessage() {}

func (x *Attr) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attr.ProtoReflect.Descriptor instead.
func (*Attr) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{2}
}

func (x *Attr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Attr) GetNameRange() *Range {
	if x != nil {
		return x.NameRange
	}
	return nil
}

func (x *Attr) GetValue() *CtyValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Attr) GetValueRange() *Range {
	if x != nil {
		return x.ValueRange
	}
	return nil
}

func (x *Attr) GetSecret() bool {
	if x != nil {
		return x.Secret
	}
	return false
}

type Block struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Header        []string               `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty"`
	HeaderRanges  []*Range               `protobuf:"bytes,2,rep,name=header_ranges,json=headerRanges,proto3" json:"header_ranges,omitempty"`
	Attributes    map[string]*Attr       `protobuf:"bytes,3,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Blocks        []*Block               `protobuf:"bytes,4,rep,name=blocks,proto3" json:"blocks,omitempty"`
	ContentsRange *Range                 `protobuf:"bytes,5,opt,name=contents_range,json=contentsRange,proto3" json:"contents_range,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Block) Reset() {
	*x = Block{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{3}
}

func (x *Block) GetHeader() []string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetHeaderRanges() []*Range {
	if x != nil {
		return x.HeaderRanges
	}
	return nil
}

func (x *Block) GetAttributes() map[string]*Attr {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Block) GetBlocks() []*Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

func (x *Block) GetContentsRange() *Range {
	if x != nil {
		return x.ContentsRange
	}
	return nil
}

type BlockSpec_NameMatcher struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Matcher:
	//
	//	*BlockSpec_NameMatcher_Exact_
	Matcher       isBlockSpec_NameMatcher_Matcher `protobuf_oneof:"matcher"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BlockSpec_NameMatcher) Reset() {
	*x = BlockSpec_NameMatcher{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockSpec_NameMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockSpec_NameMatcher) ProtoMessage() {}

func (x *BlockSpec_NameMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockSpec_NameMatcher.ProtoReflect.Descriptor instead.
func (*BlockSpec_NameMatcher) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{1, 0}
}

func (x *BlockSpec_NameMatcher) GetMatcher() isBlockSpec_NameMatcher_Matcher {
	if x != nil {
		return x.Matcher
	}
	return nil
}

func (x *BlockSpec_NameMatcher) GetExact() *BlockSpec_NameMatcher_Exact {
	if x != nil {
		if x, ok := x.Matcher.(*BlockSpec_NameMatcher_Exact_); ok {
			return x.Exact
		}
	}
	return nil
}

type isBlockSpec_NameMatcher_Matcher interface {
	isBlockSpec_NameMatcher_Matcher()
}

type BlockSpec_NameMatcher_Exact_ struct {
	Exact *BlockSpec_NameMatcher_Exact `protobuf:"bytes,1,opt,name=exact,proto3,oneof"`
}

func (*BlockSpec_NameMatcher_Exact_) isBlockSpec_NameMatcher_Matcher() {}

type BlockSpec_NameMatcher_Exact struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Matches       []string               `protobuf:"bytes,1,rep,name=matches,proto3" json:"matches,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BlockSpec_NameMatcher_Exact) Reset() {
	*x = BlockSpec_NameMatcher_Exact{}
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockSpec_NameMatcher_Exact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockSpec_NameMatcher_Exact) ProtoMessage() {}

func (x *BlockSpec_NameMatcher_Exact) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_dataspec_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockSpec_NameMatcher_Exact.ProtoReflect.Descriptor instead.
func (*BlockSpec_NameMatcher_Exact) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_dataspec_proto_rawDescGZIP(), []int{1, 0, 0}
}

func (x *BlockSpec_NameMatcher_Exact) GetMatches() []string {
	if x != nil {
		return x.Matches
	}
	return nil
}

var File_pluginapi_v1_dataspec_proto protoreflect.FileDescriptor

var file_pluginapi_v1_dataspec_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x68, 0x63, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd0, 0x03, 0x0a, 0x08,
	0x41, 0x74, 0x74, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c,
	0x12, 0x37, 0x0a, 0x0b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x63,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x6f, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x63,
	0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2d, 0x0a,
	0x06, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x74, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x6f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x3b, 0x0a, 0x0d,
	0x6d, 0x69, 0x6e, 0x5f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x74, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0c, 0x6d, 0x69, 0x6e,
	0x49, 0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x12, 0x3b, 0x0a, 0x0d, 0x6d, 0x61, 0x78,
	0x5f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x74, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x49, 0x6e, 0x63,
	0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x72,
	0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x8e,
	0x04, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x12, 0x46, 0x0a, 0x0c,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4e, 0x61, 0x6d, 0x65,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x0b, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x53, 0x70, 0x65, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64,
	0x6f, 0x63, 0x12, 0x38, 0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x70, 0x65, 0x63,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x70, 0x65, 0x63,
	0x52, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x73, 0x12, 0x35, 0x0a, 0x0a,
	0x61, 0x74, 0x74, 0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x74, 0x74, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x09, 0x61, 0x74, 0x74, 0x72, 0x53, 0x70,
	0x65, 0x63, 0x73, 0x12, 0x38, 0x0a, 0x18, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x75, 0x6e, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x55, 0x6e, 0x73, 0x70,
	0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x40, 0x0a,
	0x1c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x75, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x1a, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x1a,
	0x7e, 0x0a, 0x0b, 0x4e, 0x61, 0x6d, 0x65, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x41,
	0x0a, 0x05, 0x65, 0x78, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x2e, 0x45, 0x78, 0x61, 0x63, 0x74, 0x48, 0x00, 0x52, 0x05, 0x65, 0x78, 0x61, 0x63,
	0x74, 0x1a, 0x21, 0x0a, 0x05, 0x45, 0x78, 0x61, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x22,
	0xca, 0x01, 0x0a, 0x04, 0x41, 0x74, 0x74, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x0a,
	0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x74, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x34,
	0x0a, 0x0b, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0xda, 0x02, 0x0a,
	0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x38,
	0x0a, 0x0d, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0c, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x43, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x2b, 0x0a,
	0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x3a, 0x0a, 0x0e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x73, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x1a, 0x51, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0xb3, 0x01, 0x0a, 0x10, 0x63, 0x6f,
	0x6d, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x0d,
	0x44, 0x61, 0x74, 0x61, 0x73, 0x70, 0x65, 0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x61, 0x63,
	0x6b, 0x73, 0x74, 0x6f, 0x72, 0x6b, 0x2d, 0x69, 0x6f, 0x2f, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63,
	0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61,
	0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70,
	0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_pluginapi_v1_dataspec_proto_rawDescOnce sync.Once
	file_pluginapi_v1_dataspec_proto_rawDescData []byte
)

func file_pluginapi_v1_dataspec_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_dataspec_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_dataspec_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pluginapi_v1_dataspec_proto_rawDesc), len(file_pluginapi_v1_dataspec_proto_rawDesc)))
	})
	return file_pluginapi_v1_dataspec_proto_rawDescData
}

var file_pluginapi_v1_dataspec_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pluginapi_v1_dataspec_proto_goTypes = []any{
	(*AttrSpec)(nil),                    // 0: pluginapi.v1.AttrSpec
	(*BlockSpec)(nil),                   // 1: pluginapi.v1.BlockSpec
	(*Attr)(nil),                        // 2: pluginapi.v1.Attr
	(*Block)(nil),                       // 3: pluginapi.v1.Block
	(*BlockSpec_NameMatcher)(nil),       // 4: pluginapi.v1.BlockSpec.NameMatcher
	(*BlockSpec_NameMatcher_Exact)(nil), // 5: pluginapi.v1.BlockSpec.NameMatcher.Exact
	nil,                                 // 6: pluginapi.v1.Block.AttributesEntry
	(*CtyType)(nil),                     // 7: pluginapi.v1.CtyType
	(*CtyValue)(nil),                    // 8: pluginapi.v1.CtyValue
	(*Range)(nil),                       // 9: pluginapi.v1.Range
}
var file_pluginapi_v1_dataspec_proto_depIdxs = []int32{
	7,  // 0: pluginapi.v1.AttrSpec.type:type_name -> pluginapi.v1.CtyType
	8,  // 1: pluginapi.v1.AttrSpec.default_val:type_name -> pluginapi.v1.CtyValue
	8,  // 2: pluginapi.v1.AttrSpec.example_val:type_name -> pluginapi.v1.CtyValue
	8,  // 3: pluginapi.v1.AttrSpec.one_of:type_name -> pluginapi.v1.CtyValue
	8,  // 4: pluginapi.v1.AttrSpec.min_inclusive:type_name -> pluginapi.v1.CtyValue
	8,  // 5: pluginapi.v1.AttrSpec.max_inclusive:type_name -> pluginapi.v1.CtyValue
	4,  // 6: pluginapi.v1.BlockSpec.headers_spec:type_name -> pluginapi.v1.BlockSpec.NameMatcher
	1,  // 7: pluginapi.v1.BlockSpec.block_specs:type_name -> pluginapi.v1.BlockSpec
	0,  // 8: pluginapi.v1.BlockSpec.attr_specs:type_name -> pluginapi.v1.AttrSpec
	9,  // 9: pluginapi.v1.Attr.name_range:type_name -> pluginapi.v1.Range
	8,  // 10: pluginapi.v1.Attr.value:type_name -> pluginapi.v1.CtyValue
	9,  // 11: pluginapi.v1.Attr.value_range:type_name -> pluginapi.v1.Range
	9,  // 12: pluginapi.v1.Block.header_ranges:type_name -> pluginapi.v1.Range
	6,  // 13: pluginapi.v1.Block.attributes:type_name -> pluginapi.v1.Block.AttributesEntry
	3,  // 14: pluginapi.v1.Block.blocks:type_name -> pluginapi.v1.Block
	9,  // 15: pluginapi.v1.Block.contents_range:type_name -> pluginapi.v1.Range
	5,  // 16: pluginapi.v1.BlockSpec.NameMatcher.exact:type_name -> pluginapi.v1.BlockSpec.NameMatcher.Exact
	2,  // 17: pluginapi.v1.Block.AttributesEntry.value:type_name -> pluginapi.v1.Attr
	18, // [18:18] is the sub-list for method output_type
	18, // [18:18] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_dataspec_proto_init() }
func file_pluginapi_v1_dataspec_proto_init() {
	if File_pluginapi_v1_dataspec_proto != nil {
		return
	}
	file_pluginapi_v1_cty_proto_init()
	file_pluginapi_v1_hcl_proto_init()
	file_pluginapi_v1_dataspec_proto_msgTypes[4].OneofWrappers = []any{
		(*BlockSpec_NameMatcher_Exact_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pluginapi_v1_dataspec_proto_rawDesc), len(file_pluginapi_v1_dataspec_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_dataspec_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_dataspec_proto_depIdxs,
		MessageInfos:      file_pluginapi_v1_dataspec_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_dataspec_proto = out.File
	file_pluginapi_v1_dataspec_proto_goTypes = nil
	file_pluginapi_v1_dataspec_proto_depIdxs = nil
}
