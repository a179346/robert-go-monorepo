// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.28.2
// source: pb/calculator/calculator.proto

package calculator

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

type DivideRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dividend      float64                `protobuf:"fixed64,1,opt,name=dividend,proto3" json:"dividend,omitempty"`
	Divisor       float64                `protobuf:"fixed64,2,opt,name=divisor,proto3" json:"divisor,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DivideRequest) Reset() {
	*x = DivideRequest{}
	mi := &file_pb_calculator_calculator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DivideRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DivideRequest) ProtoMessage() {}

func (x *DivideRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_calculator_calculator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DivideRequest.ProtoReflect.Descriptor instead.
func (*DivideRequest) Descriptor() ([]byte, []int) {
	return file_pb_calculator_calculator_proto_rawDescGZIP(), []int{0}
}

func (x *DivideRequest) GetDividend() float64 {
	if x != nil {
		return x.Dividend
	}
	return 0
}

func (x *DivideRequest) GetDivisor() float64 {
	if x != nil {
		return x.Divisor
	}
	return 0
}

type DivideResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         float64                `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DivideResponse) Reset() {
	*x = DivideResponse{}
	mi := &file_pb_calculator_calculator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DivideResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DivideResponse) ProtoMessage() {}

func (x *DivideResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_calculator_calculator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DivideResponse.ProtoReflect.Descriptor instead.
func (*DivideResponse) Descriptor() ([]byte, []int) {
	return file_pb_calculator_calculator_proto_rawDescGZIP(), []int{1}
}

func (x *DivideResponse) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_pb_calculator_calculator_proto protoreflect.FileDescriptor

var file_pb_calculator_calculator_proto_rawDesc = string([]byte{
	0x0a, 0x1e, 0x70, 0x62, 0x2f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2f,
	0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x45, 0x0a, 0x0d,
	0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x08, 0x64, 0x69, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x76,
	0x69, 0x73, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x64, 0x69, 0x76, 0x69,
	0x73, 0x6f, 0x72, 0x22, 0x26, 0x0a, 0x0e, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x4f, 0x0a, 0x0a, 0x43,
	0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x41, 0x0a, 0x06, 0x44, 0x69, 0x76,
	0x69, 0x64, 0x65, 0x12, 0x19, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x69, 0x76, 0x69,
	0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x35, 0x5a, 0x33,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x31, 0x37, 0x39, 0x33,
	0x34, 0x36, 0x2f, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x74, 0x2d, 0x67, 0x6f, 0x2d, 0x6d, 0x6f, 0x6e,
	0x6f, 0x72, 0x65, 0x70, 0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_pb_calculator_calculator_proto_rawDescOnce sync.Once
	file_pb_calculator_calculator_proto_rawDescData []byte
)

func file_pb_calculator_calculator_proto_rawDescGZIP() []byte {
	file_pb_calculator_calculator_proto_rawDescOnce.Do(func() {
		file_pb_calculator_calculator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pb_calculator_calculator_proto_rawDesc), len(file_pb_calculator_calculator_proto_rawDesc)))
	})
	return file_pb_calculator_calculator_proto_rawDescData
}

var file_pb_calculator_calculator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_calculator_calculator_proto_goTypes = []any{
	(*DivideRequest)(nil),  // 0: calculator.DivideRequest
	(*DivideResponse)(nil), // 1: calculator.DivideResponse
}
var file_pb_calculator_calculator_proto_depIdxs = []int32{
	0, // 0: calculator.Calculator.Divide:input_type -> calculator.DivideRequest
	1, // 1: calculator.Calculator.Divide:output_type -> calculator.DivideResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_calculator_calculator_proto_init() }
func file_pb_calculator_calculator_proto_init() {
	if File_pb_calculator_calculator_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pb_calculator_calculator_proto_rawDesc), len(file_pb_calculator_calculator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_calculator_calculator_proto_goTypes,
		DependencyIndexes: file_pb_calculator_calculator_proto_depIdxs,
		MessageInfos:      file_pb_calculator_calculator_proto_msgTypes,
	}.Build()
	File_pb_calculator_calculator_proto = out.File
	file_pb_calculator_calculator_proto_goTypes = nil
	file_pb_calculator_calculator_proto_depIdxs = nil
}
