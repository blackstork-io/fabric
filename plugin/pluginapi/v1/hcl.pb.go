// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: pluginapi/v1/hcl.proto

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

type Pos struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Line          int64                  `protobuf:"varint,1,opt,name=line,proto3" json:"line,omitempty"`
	Column        int64                  `protobuf:"varint,2,opt,name=column,proto3" json:"column,omitempty"`
	Byte          int64                  `protobuf:"varint,3,opt,name=byte,proto3" json:"byte,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pos) Reset() {
	*x = Pos{}
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pos) ProtoMessage() {}

func (x *Pos) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pos.ProtoReflect.Descriptor instead.
func (*Pos) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_hcl_proto_rawDescGZIP(), []int{0}
}

func (x *Pos) GetLine() int64 {
	if x != nil {
		return x.Line
	}
	return 0
}

func (x *Pos) GetColumn() int64 {
	if x != nil {
		return x.Column
	}
	return 0
}

func (x *Pos) GetByte() int64 {
	if x != nil {
		return x.Byte
	}
	return 0
}

type Range struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Start         *Pos                   `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End           *Pos                   `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Range) Reset() {
	*x = Range{}
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Range) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Range) ProtoMessage() {}

func (x *Range) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Range.ProtoReflect.Descriptor instead.
func (*Range) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_hcl_proto_rawDescGZIP(), []int{1}
}

func (x *Range) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Range) GetStart() *Pos {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *Range) GetEnd() *Pos {
	if x != nil {
		return x.End
	}
	return nil
}

type Diagnostic struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Severity      int64                  `protobuf:"varint,1,opt,name=severity,proto3" json:"severity,omitempty"`
	Summary       string                 `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	Detail        string                 `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	Subject       *Range                 `protobuf:"bytes,4,opt,name=subject,proto3" json:"subject,omitempty"`
	Context       *Range                 `protobuf:"bytes,5,opt,name=context,proto3" json:"context,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Diagnostic) Reset() {
	*x = Diagnostic{}
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Diagnostic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Diagnostic) ProtoMessage() {}

func (x *Diagnostic) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_hcl_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Diagnostic.ProtoReflect.Descriptor instead.
func (*Diagnostic) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_hcl_proto_rawDescGZIP(), []int{2}
}

func (x *Diagnostic) GetSeverity() int64 {
	if x != nil {
		return x.Severity
	}
	return 0
}

func (x *Diagnostic) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Diagnostic) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *Diagnostic) GetSubject() *Range {
	if x != nil {
		return x.Subject
	}
	return nil
}

func (x *Diagnostic) GetContext() *Range {
	if x != nil {
		return x.Context
	}
	return nil
}

var File_pluginapi_v1_hcl_proto protoreflect.FileDescriptor

var file_pluginapi_v1_hcl_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x68,
	0x63, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x22, 0x45, 0x0a, 0x03, 0x50, 0x6f, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6c, 0x69, 0x6e,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x79, 0x74,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x62, 0x79, 0x74, 0x65, 0x22, 0x71, 0x0a,
	0x05, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x6f, 0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x23, 0x0a, 0x03, 0x65,
	0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x52, 0x03, 0x65, 0x6e, 0x64,
	0x22, 0xb8, 0x01, 0x0a, 0x0a, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x12,
	0x1a, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x2d, 0x0a,
	0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x2d, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x42, 0xae, 0x01, 0x0a, 0x10,
	0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x42, 0x08, 0x48, 0x63, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b, 0x73, 0x74,
	0x6f, 0x72, 0x6b, 0x2d, 0x69, 0x6f, 0x2f, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x2f, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x3b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x50, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x0c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x18, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pluginapi_v1_hcl_proto_rawDescOnce sync.Once
	file_pluginapi_v1_hcl_proto_rawDescData = file_pluginapi_v1_hcl_proto_rawDesc
)

func file_pluginapi_v1_hcl_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_hcl_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_hcl_proto_rawDescData = protoimpl.X.CompressGZIP(file_pluginapi_v1_hcl_proto_rawDescData)
	})
	return file_pluginapi_v1_hcl_proto_rawDescData
}

var file_pluginapi_v1_hcl_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pluginapi_v1_hcl_proto_goTypes = []any{
	(*Pos)(nil),        // 0: pluginapi.v1.Pos
	(*Range)(nil),      // 1: pluginapi.v1.Range
	(*Diagnostic)(nil), // 2: pluginapi.v1.Diagnostic
}
var file_pluginapi_v1_hcl_proto_depIdxs = []int32{
	0, // 0: pluginapi.v1.Range.start:type_name -> pluginapi.v1.Pos
	0, // 1: pluginapi.v1.Range.end:type_name -> pluginapi.v1.Pos
	1, // 2: pluginapi.v1.Diagnostic.subject:type_name -> pluginapi.v1.Range
	1, // 3: pluginapi.v1.Diagnostic.context:type_name -> pluginapi.v1.Range
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_hcl_proto_init() }
func file_pluginapi_v1_hcl_proto_init() {
	if File_pluginapi_v1_hcl_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pluginapi_v1_hcl_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_hcl_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_hcl_proto_depIdxs,
		MessageInfos:      file_pluginapi_v1_hcl_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_hcl_proto = out.File
	file_pluginapi_v1_hcl_proto_rawDesc = nil
	file_pluginapi_v1_hcl_proto_goTypes = nil
	file_pluginapi_v1_hcl_proto_depIdxs = nil
}
