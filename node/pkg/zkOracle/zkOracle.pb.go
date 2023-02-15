// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: zkOracle.proto

package zkOracle

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

type SendVoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index     uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Request   uint64 `protobuf:"varint,2,opt,name=request,proto3" json:"request,omitempty"`
	BlockHash []byte `protobuf:"bytes,3,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *SendVoteRequest) Reset() {
	*x = SendVoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zkOracle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendVoteRequest) ProtoMessage() {}

func (x *SendVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_zkOracle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendVoteRequest.ProtoReflect.Descriptor instead.
func (*SendVoteRequest) Descriptor() ([]byte, []int) {
	return file_zkOracle_proto_rawDescGZIP(), []int{0}
}

func (x *SendVoteRequest) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *SendVoteRequest) GetRequest() uint64 {
	if x != nil {
		return x.Request
	}
	return 0
}

func (x *SendVoteRequest) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *SendVoteRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type SendVoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendVoteResponse) Reset() {
	*x = SendVoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zkOracle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendVoteResponse) ProtoMessage() {}

func (x *SendVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_zkOracle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendVoteResponse.ProtoReflect.Descriptor instead.
func (*SendVoteResponse) Descriptor() ([]byte, []int) {
	return file_zkOracle_proto_rawDescGZIP(), []int{1}
}

var File_zkOracle_proto protoreflect.FileDescriptor

var file_zkOracle_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x7a, 0x6b, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x7a, 0x6b, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x22, 0x7d, 0x0a, 0x0f, 0x53, 0x65,
	0x6e, 0x64, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x53, 0x65, 0x6e,
	0x64, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x4f, 0x0a,
	0x0a, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x41, 0x0a, 0x08, 0x53,
	0x65, 0x6e, 0x64, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x7a, 0x6b, 0x4f, 0x72, 0x61, 0x63,
	0x6c, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x7a, 0x6b, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x2e, 0x53, 0x65,
	0x6e, 0x64, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x3b, 0x7a, 0x6b, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_zkOracle_proto_rawDescOnce sync.Once
	file_zkOracle_proto_rawDescData = file_zkOracle_proto_rawDesc
)

func file_zkOracle_proto_rawDescGZIP() []byte {
	file_zkOracle_proto_rawDescOnce.Do(func() {
		file_zkOracle_proto_rawDescData = protoimpl.X.CompressGZIP(file_zkOracle_proto_rawDescData)
	})
	return file_zkOracle_proto_rawDescData
}

var file_zkOracle_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_zkOracle_proto_goTypes = []interface{}{
	(*SendVoteRequest)(nil),  // 0: zkOracle.SendVoteRequest
	(*SendVoteResponse)(nil), // 1: zkOracle.SendVoteResponse
}
var file_zkOracle_proto_depIdxs = []int32{
	0, // 0: zkOracle.OracleNode.SendVote:input_type -> zkOracle.SendVoteRequest
	1, // 1: zkOracle.OracleNode.SendVote:output_type -> zkOracle.SendVoteResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_zkOracle_proto_init() }
func file_zkOracle_proto_init() {
	if File_zkOracle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zkOracle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendVoteRequest); i {
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
		file_zkOracle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendVoteResponse); i {
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
			RawDescriptor: file_zkOracle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_zkOracle_proto_goTypes,
		DependencyIndexes: file_zkOracle_proto_depIdxs,
		MessageInfos:      file_zkOracle_proto_msgTypes,
	}.Build()
	File_zkOracle_proto = out.File
	file_zkOracle_proto_rawDesc = nil
	file_zkOracle_proto_goTypes = nil
	file_zkOracle_proto_depIdxs = nil
}
