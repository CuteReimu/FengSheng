// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.11.4
// source: errcode.proto

package protos

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

type ErrorCode int32

const (
	ErrorCode_client_version_not_match ErrorCode = 0 // 客户端版本号不匹配，int_params[0]为服务器版本号
	ErrorCode_no_more_room             ErrorCode = 1 // 没有更多的房间了
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "client_version_not_match",
		1: "no_more_room",
	}
	ErrorCode_value = map[string]int32{
		"client_version_not_match": 0,
		"no_more_room":             1,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_errcode_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_errcode_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_errcode_proto_rawDescGZIP(), []int{0}
}

type ErrorCodeToc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code      ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=ErrorCode" json:"code,omitempty"`
	IntParams []int64   `protobuf:"varint,2,rep,packed,name=int_params,json=intParams,proto3" json:"int_params,omitempty"`
	StrParams []string  `protobuf:"bytes,3,rep,name=str_params,json=strParams,proto3" json:"str_params,omitempty"`
}

func (x *ErrorCodeToc) Reset() {
	*x = ErrorCodeToc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errcode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorCodeToc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorCodeToc) ProtoMessage() {}

func (x *ErrorCodeToc) ProtoReflect() protoreflect.Message {
	mi := &file_errcode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorCodeToc.ProtoReflect.Descriptor instead.
func (*ErrorCodeToc) Descriptor() ([]byte, []int) {
	return file_errcode_proto_rawDescGZIP(), []int{0}
}

func (x *ErrorCodeToc) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_client_version_not_match
}

func (x *ErrorCodeToc) GetIntParams() []int64 {
	if x != nil {
		return x.IntParams
	}
	return nil
}

func (x *ErrorCodeToc) GetStrParams() []string {
	if x != nil {
		return x.StrParams
	}
	return nil
}

var File_errcode_proto protoreflect.FileDescriptor

var file_errcode_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x6f, 0x0a, 0x0e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x6f,
	0x63, 0x12, 0x1f, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0b, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x72, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x2a, 0x3c, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1c,
	0x0a, 0x18, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x5f, 0x6e, 0x6f, 0x74, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c,
	0x6e, 0x6f, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x6d, 0x10, 0x01, 0x42, 0x10,
	0x5a, 0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errcode_proto_rawDescOnce sync.Once
	file_errcode_proto_rawDescData = file_errcode_proto_rawDesc
)

func file_errcode_proto_rawDescGZIP() []byte {
	file_errcode_proto_rawDescOnce.Do(func() {
		file_errcode_proto_rawDescData = protoimpl.X.CompressGZIP(file_errcode_proto_rawDescData)
	})
	return file_errcode_proto_rawDescData
}

var file_errcode_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errcode_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_errcode_proto_goTypes = []interface{}{
	(ErrorCode)(0),       // 0: error_code
	(*ErrorCodeToc)(nil), // 1: error_code_toc
}
var file_errcode_proto_depIdxs = []int32{
	0, // 0: error_code_toc.code:type_name -> error_code
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_errcode_proto_init() }
func file_errcode_proto_init() {
	if File_errcode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_errcode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorCodeToc); i {
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
			RawDescriptor: file_errcode_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errcode_proto_goTypes,
		DependencyIndexes: file_errcode_proto_depIdxs,
		EnumInfos:         file_errcode_proto_enumTypes,
		MessageInfos:      file_errcode_proto_msgTypes,
	}.Build()
	File_errcode_proto = out.File
	file_errcode_proto_rawDesc = nil
	file_errcode_proto_goTypes = nil
	file_errcode_proto_depIdxs = nil
}