// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: pluginapi/v1/schema.proto

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

type InvocationOrder int32

const (
	InvocationOrder_INVOCATION_ORDER_UNSPECIFIED InvocationOrder = 0
	InvocationOrder_INVOCATION_ORDER_BEGIN       InvocationOrder = 2
	InvocationOrder_INVOCATION_ORDER_END         InvocationOrder = 3
)

// Enum value maps for InvocationOrder.
var (
	InvocationOrder_name = map[int32]string{
		0: "INVOCATION_ORDER_UNSPECIFIED",
		2: "INVOCATION_ORDER_BEGIN",
		3: "INVOCATION_ORDER_END",
	}
	InvocationOrder_value = map[string]int32{
		"INVOCATION_ORDER_UNSPECIFIED": 0,
		"INVOCATION_ORDER_BEGIN":       2,
		"INVOCATION_ORDER_END":         3,
	}
)

func (x InvocationOrder) Enum() *InvocationOrder {
	p := new(InvocationOrder)
	*p = x
	return p
}

func (x InvocationOrder) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InvocationOrder) Descriptor() protoreflect.EnumDescriptor {
	return file_pluginapi_v1_schema_proto_enumTypes[0].Descriptor()
}

func (InvocationOrder) Type() protoreflect.EnumType {
	return &file_pluginapi_v1_schema_proto_enumTypes[0]
}

func (x InvocationOrder) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InvocationOrder.Descriptor instead.
func (InvocationOrder) EnumDescriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{0}
}

type OutputFormat int32

const (
	OutputFormat_OUTPUT_FORMAT_UNSPECIFIED OutputFormat = 0
	OutputFormat_OUTPUT_FORMAT_MD          OutputFormat = 1
	OutputFormat_OUTPUT_FORMAT_HTML        OutputFormat = 2
	OutputFormat_OUTPUT_FORMAT_PDF         OutputFormat = 3
)

// Enum value maps for OutputFormat.
var (
	OutputFormat_name = map[int32]string{
		0: "OUTPUT_FORMAT_UNSPECIFIED",
		1: "OUTPUT_FORMAT_MD",
		2: "OUTPUT_FORMAT_HTML",
		3: "OUTPUT_FORMAT_PDF",
	}
	OutputFormat_value = map[string]int32{
		"OUTPUT_FORMAT_UNSPECIFIED": 0,
		"OUTPUT_FORMAT_MD":          1,
		"OUTPUT_FORMAT_HTML":        2,
		"OUTPUT_FORMAT_PDF":         3,
	}
)

func (x OutputFormat) Enum() *OutputFormat {
	p := new(OutputFormat)
	*p = x
	return p
}

func (x OutputFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OutputFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_pluginapi_v1_schema_proto_enumTypes[1].Descriptor()
}

func (OutputFormat) Type() protoreflect.EnumType {
	return &file_pluginapi_v1_schema_proto_enumTypes[1]
}

func (x OutputFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OutputFormat.Descriptor instead.
func (OutputFormat) EnumDescriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{1}
}

type Schema struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Name    string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Plugin components
	DataSources      map[string]*DataSourceSchema      `protobuf:"bytes,3,rep,name=data_sources,json=dataSources,proto3" json:"data_sources,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ContentProviders map[string]*ContentProviderSchema `protobuf:"bytes,4,rep,name=content_providers,json=contentProviders,proto3" json:"content_providers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Publishers       map[string]*PublisherSchema       `protobuf:"bytes,7,rep,name=publishers,proto3" json:"publishers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Doc              string                            `protobuf:"bytes,5,opt,name=doc,proto3" json:"doc,omitempty"`
	Tags             []string                          `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Schema) Reset() {
	*x = Schema{}
	mi := &file_pluginapi_v1_schema_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schema) ProtoMessage() {}

func (x *Schema) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_schema_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schema.ProtoReflect.Descriptor instead.
func (*Schema) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{0}
}

func (x *Schema) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Schema) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Schema) GetDataSources() map[string]*DataSourceSchema {
	if x != nil {
		return x.DataSources
	}
	return nil
}

func (x *Schema) GetContentProviders() map[string]*ContentProviderSchema {
	if x != nil {
		return x.ContentProviders
	}
	return nil
}

func (x *Schema) GetPublishers() map[string]*PublisherSchema {
	if x != nil {
		return x.Publishers
	}
	return nil
}

func (x *Schema) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *Schema) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type DataSourceSchema struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Args          *BlockSpec             `protobuf:"bytes,3,opt,name=args,proto3" json:"args,omitempty"`
	Config        *BlockSpec             `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
	Doc           string                 `protobuf:"bytes,5,opt,name=doc,proto3" json:"doc,omitempty"`
	Tags          []string               `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataSourceSchema) Reset() {
	*x = DataSourceSchema{}
	mi := &file_pluginapi_v1_schema_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataSourceSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSourceSchema) ProtoMessage() {}

func (x *DataSourceSchema) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_schema_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSourceSchema.ProtoReflect.Descriptor instead.
func (*DataSourceSchema) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{1}
}

func (x *DataSourceSchema) GetArgs() *BlockSpec {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *DataSourceSchema) GetConfig() *BlockSpec {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *DataSourceSchema) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *DataSourceSchema) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type ContentProviderSchema struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Args            *BlockSpec             `protobuf:"bytes,4,opt,name=args,proto3" json:"args,omitempty"`
	Config          *BlockSpec             `protobuf:"bytes,5,opt,name=config,proto3" json:"config,omitempty"`
	InvocationOrder InvocationOrder        `protobuf:"varint,3,opt,name=invocation_order,json=invocationOrder,proto3,enum=pluginapi.v1.InvocationOrder" json:"invocation_order,omitempty"`
	Doc             string                 `protobuf:"bytes,6,opt,name=doc,proto3" json:"doc,omitempty"`
	Tags            []string               `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ContentProviderSchema) Reset() {
	*x = ContentProviderSchema{}
	mi := &file_pluginapi_v1_schema_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentProviderSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentProviderSchema) ProtoMessage() {}

func (x *ContentProviderSchema) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_schema_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentProviderSchema.ProtoReflect.Descriptor instead.
func (*ContentProviderSchema) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{2}
}

func (x *ContentProviderSchema) GetArgs() *BlockSpec {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *ContentProviderSchema) GetConfig() *BlockSpec {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *ContentProviderSchema) GetInvocationOrder() InvocationOrder {
	if x != nil {
		return x.InvocationOrder
	}
	return InvocationOrder_INVOCATION_ORDER_UNSPECIFIED
}

func (x *ContentProviderSchema) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *ContentProviderSchema) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type PublisherSchema struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Args           *BlockSpec             `protobuf:"bytes,1,opt,name=args,proto3" json:"args,omitempty"`
	Config         *BlockSpec             `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	Doc            string                 `protobuf:"bytes,3,opt,name=doc,proto3" json:"doc,omitempty"`
	Tags           []string               `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	AllowedFormats []OutputFormat         `protobuf:"varint,5,rep,packed,name=allowed_formats,json=allowedFormats,proto3,enum=pluginapi.v1.OutputFormat" json:"allowed_formats,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PublisherSchema) Reset() {
	*x = PublisherSchema{}
	mi := &file_pluginapi_v1_schema_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PublisherSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublisherSchema) ProtoMessage() {}

func (x *PublisherSchema) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_v1_schema_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublisherSchema.ProtoReflect.Descriptor instead.
func (*PublisherSchema) Descriptor() ([]byte, []int) {
	return file_pluginapi_v1_schema_proto_rawDescGZIP(), []int{3}
}

func (x *PublisherSchema) GetArgs() *BlockSpec {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *PublisherSchema) GetConfig() *BlockSpec {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *PublisherSchema) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

func (x *PublisherSchema) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *PublisherSchema) GetAllowedFormats() []OutputFormat {
	if x != nil {
		return x.AllowedFormats
	}
	return nil
}

var File_pluginapi_v1_schema_proto protoreflect.FileDescriptor

const file_pluginapi_v1_schema_proto_rawDesc = "" +
	"\n" +
	"\x19pluginapi/v1/schema.proto\x12\fpluginapi.v1\x1a\x1bpluginapi/v1/dataspec.proto\"\xed\x04\n" +
	"\x06Schema\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x18\n" +
	"\aversion\x18\x02 \x01(\tR\aversion\x12H\n" +
	"\fdata_sources\x18\x03 \x03(\v2%.pluginapi.v1.Schema.DataSourcesEntryR\vdataSources\x12W\n" +
	"\x11content_providers\x18\x04 \x03(\v2*.pluginapi.v1.Schema.ContentProvidersEntryR\x10contentProviders\x12D\n" +
	"\n" +
	"publishers\x18\a \x03(\v2$.pluginapi.v1.Schema.PublishersEntryR\n" +
	"publishers\x12\x10\n" +
	"\x03doc\x18\x05 \x01(\tR\x03doc\x12\x12\n" +
	"\x04tags\x18\x06 \x03(\tR\x04tags\x1a^\n" +
	"\x10DataSourcesEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x124\n" +
	"\x05value\x18\x02 \x01(\v2\x1e.pluginapi.v1.DataSourceSchemaR\x05value:\x028\x01\x1ah\n" +
	"\x15ContentProvidersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x129\n" +
	"\x05value\x18\x02 \x01(\v2#.pluginapi.v1.ContentProviderSchemaR\x05value:\x028\x01\x1a\\\n" +
	"\x0fPublishersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x123\n" +
	"\x05value\x18\x02 \x01(\v2\x1d.pluginapi.v1.PublisherSchemaR\x05value:\x028\x01\"\x96\x01\n" +
	"\x10DataSourceSchema\x12+\n" +
	"\x04args\x18\x03 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x04args\x12/\n" +
	"\x06config\x18\x04 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x06config\x12\x10\n" +
	"\x03doc\x18\x05 \x01(\tR\x03doc\x12\x12\n" +
	"\x04tags\x18\x06 \x03(\tR\x04tags\"\xe5\x01\n" +
	"\x15ContentProviderSchema\x12+\n" +
	"\x04args\x18\x04 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x04args\x12/\n" +
	"\x06config\x18\x05 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x06config\x12H\n" +
	"\x10invocation_order\x18\x03 \x01(\x0e2\x1d.pluginapi.v1.InvocationOrderR\x0finvocationOrder\x12\x10\n" +
	"\x03doc\x18\x06 \x01(\tR\x03doc\x12\x12\n" +
	"\x04tags\x18\a \x03(\tR\x04tags\"\xda\x01\n" +
	"\x0fPublisherSchema\x12+\n" +
	"\x04args\x18\x01 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x04args\x12/\n" +
	"\x06config\x18\x02 \x01(\v2\x17.pluginapi.v1.BlockSpecR\x06config\x12\x10\n" +
	"\x03doc\x18\x03 \x01(\tR\x03doc\x12\x12\n" +
	"\x04tags\x18\x04 \x03(\tR\x04tags\x12C\n" +
	"\x0fallowed_formats\x18\x05 \x03(\x0e2\x1a.pluginapi.v1.OutputFormatR\x0eallowedFormats*i\n" +
	"\x0fInvocationOrder\x12 \n" +
	"\x1cINVOCATION_ORDER_UNSPECIFIED\x10\x00\x12\x1a\n" +
	"\x16INVOCATION_ORDER_BEGIN\x10\x02\x12\x18\n" +
	"\x14INVOCATION_ORDER_END\x10\x03*r\n" +
	"\fOutputFormat\x12\x1d\n" +
	"\x19OUTPUT_FORMAT_UNSPECIFIED\x10\x00\x12\x14\n" +
	"\x10OUTPUT_FORMAT_MD\x10\x01\x12\x16\n" +
	"\x12OUTPUT_FORMAT_HTML\x10\x02\x12\x15\n" +
	"\x11OUTPUT_FORMAT_PDF\x10\x03B\xb1\x01\n" +
	"\x10com.pluginapi.v1B\vSchemaProtoP\x01Z?github.com/blackstork-io/fabric/plugin/pluginapi/v1;pluginapiv1\xa2\x02\x03PXX\xaa\x02\fPluginapi.V1\xca\x02\fPluginapi\\V1\xe2\x02\x18Pluginapi\\V1\\GPBMetadata\xea\x02\rPluginapi::V1b\x06proto3"

var (
	file_pluginapi_v1_schema_proto_rawDescOnce sync.Once
	file_pluginapi_v1_schema_proto_rawDescData []byte
)

func file_pluginapi_v1_schema_proto_rawDescGZIP() []byte {
	file_pluginapi_v1_schema_proto_rawDescOnce.Do(func() {
		file_pluginapi_v1_schema_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pluginapi_v1_schema_proto_rawDesc), len(file_pluginapi_v1_schema_proto_rawDesc)))
	})
	return file_pluginapi_v1_schema_proto_rawDescData
}

var file_pluginapi_v1_schema_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_pluginapi_v1_schema_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pluginapi_v1_schema_proto_goTypes = []any{
	(InvocationOrder)(0),          // 0: pluginapi.v1.InvocationOrder
	(OutputFormat)(0),             // 1: pluginapi.v1.OutputFormat
	(*Schema)(nil),                // 2: pluginapi.v1.Schema
	(*DataSourceSchema)(nil),      // 3: pluginapi.v1.DataSourceSchema
	(*ContentProviderSchema)(nil), // 4: pluginapi.v1.ContentProviderSchema
	(*PublisherSchema)(nil),       // 5: pluginapi.v1.PublisherSchema
	nil,                           // 6: pluginapi.v1.Schema.DataSourcesEntry
	nil,                           // 7: pluginapi.v1.Schema.ContentProvidersEntry
	nil,                           // 8: pluginapi.v1.Schema.PublishersEntry
	(*BlockSpec)(nil),             // 9: pluginapi.v1.BlockSpec
}
var file_pluginapi_v1_schema_proto_depIdxs = []int32{
	6,  // 0: pluginapi.v1.Schema.data_sources:type_name -> pluginapi.v1.Schema.DataSourcesEntry
	7,  // 1: pluginapi.v1.Schema.content_providers:type_name -> pluginapi.v1.Schema.ContentProvidersEntry
	8,  // 2: pluginapi.v1.Schema.publishers:type_name -> pluginapi.v1.Schema.PublishersEntry
	9,  // 3: pluginapi.v1.DataSourceSchema.args:type_name -> pluginapi.v1.BlockSpec
	9,  // 4: pluginapi.v1.DataSourceSchema.config:type_name -> pluginapi.v1.BlockSpec
	9,  // 5: pluginapi.v1.ContentProviderSchema.args:type_name -> pluginapi.v1.BlockSpec
	9,  // 6: pluginapi.v1.ContentProviderSchema.config:type_name -> pluginapi.v1.BlockSpec
	0,  // 7: pluginapi.v1.ContentProviderSchema.invocation_order:type_name -> pluginapi.v1.InvocationOrder
	9,  // 8: pluginapi.v1.PublisherSchema.args:type_name -> pluginapi.v1.BlockSpec
	9,  // 9: pluginapi.v1.PublisherSchema.config:type_name -> pluginapi.v1.BlockSpec
	1,  // 10: pluginapi.v1.PublisherSchema.allowed_formats:type_name -> pluginapi.v1.OutputFormat
	3,  // 11: pluginapi.v1.Schema.DataSourcesEntry.value:type_name -> pluginapi.v1.DataSourceSchema
	4,  // 12: pluginapi.v1.Schema.ContentProvidersEntry.value:type_name -> pluginapi.v1.ContentProviderSchema
	5,  // 13: pluginapi.v1.Schema.PublishersEntry.value:type_name -> pluginapi.v1.PublisherSchema
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_pluginapi_v1_schema_proto_init() }
func file_pluginapi_v1_schema_proto_init() {
	if File_pluginapi_v1_schema_proto != nil {
		return
	}
	file_pluginapi_v1_dataspec_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pluginapi_v1_schema_proto_rawDesc), len(file_pluginapi_v1_schema_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pluginapi_v1_schema_proto_goTypes,
		DependencyIndexes: file_pluginapi_v1_schema_proto_depIdxs,
		EnumInfos:         file_pluginapi_v1_schema_proto_enumTypes,
		MessageInfos:      file_pluginapi_v1_schema_proto_msgTypes,
	}.Build()
	File_pluginapi_v1_schema_proto = out.File
	file_pluginapi_v1_schema_proto_goTypes = nil
	file_pluginapi_v1_schema_proto_depIdxs = nil
}
