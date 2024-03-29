// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: node_messages.proto

package main

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ContentChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *ContentChunk) Reset() {
	*x = ContentChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContentChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentChunk) ProtoMessage() {}

func (x *ContentChunk) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentChunk.ProtoReflect.Descriptor instead.
func (*ContentChunk) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{0}
}

func (x *ContentChunk) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

type ContentID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ContentID) Reset() {
	*x = ContentID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContentID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentID) ProtoMessage() {}

func (x *ContentID) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentID.ProtoReflect.Descriptor instead.
func (*ContentID) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{1}
}

func (x *ContentID) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type BlockPosition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ViewId         uint64 `protobuf:"varint,1,opt,name=viewId,proto3" json:"viewId,omitempty"`
	LatestSequence uint64 `protobuf:"varint,2,opt,name=latestSequence,proto3" json:"latestSequence,omitempty"`
}

func (x *BlockPosition) Reset() {
	*x = BlockPosition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockPosition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockPosition) ProtoMessage() {}

func (x *BlockPosition) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockPosition.ProtoReflect.Descriptor instead.
func (*BlockPosition) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{2}
}

func (x *BlockPosition) GetViewId() uint64 {
	if x != nil {
		return x.ViewId
	}
	return 0
}

func (x *BlockPosition) GetLatestSequence() uint64 {
	if x != nil {
		return x.LatestSequence
	}
	return 0
}

type BlockRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata []byte   `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"` //opaque view metadata
	Batch    [][]byte `protobuf:"bytes,2,rep,name=batch,proto3" json:"batch,omitempty"`
	PrevHash []byte   `protobuf:"bytes,3,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
}

func (x *BlockRecord) Reset() {
	*x = BlockRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockRecord) ProtoMessage() {}

func (x *BlockRecord) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockRecord.ProtoReflect.Descriptor instead.
func (*BlockRecord) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{3}
}

func (x *BlockRecord) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *BlockRecord) GetBatch() [][]byte {
	if x != nil {
		return x.Batch
	}
	return nil
}

func (x *BlockRecord) GetPrevHash() []byte {
	if x != nil {
		return x.PrevHash
	}
	return nil
}

type FwdMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender  uint64 `protobuf:"varint,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *FwdMessage) Reset() {
	*x = FwdMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FwdMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FwdMessage) ProtoMessage() {}

func (x *FwdMessage) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FwdMessage.ProtoReflect.Descriptor instead.
func (*FwdMessage) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{4}
}

func (x *FwdMessage) GetSender() uint64 {
	if x != nil {
		return x.Sender
	}
	return 0
}

func (x *FwdMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Consensus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node    uint64     `protobuf:"varint,1,opt,name=node,proto3" json:"node,omitempty"`
	Message *anypb.Any `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Consensus) Reset() {
	*x = Consensus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Consensus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Consensus) ProtoMessage() {}

func (x *Consensus) ProtoReflect() protoreflect.Message {
	mi := &file_node_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Consensus.ProtoReflect.Descriptor instead.
func (*Consensus) Descriptor() ([]byte, []int) {
	return file_node_messages_proto_rawDescGZIP(), []int{5}
}

func (x *Consensus) GetNode() uint64 {
	if x != nil {
		return x.Node
	}
	return 0
}

func (x *Consensus) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_node_messages_proto protoreflect.FileDescriptor

var file_node_messages_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x6f, 0x1a, 0x19,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x24, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x1b, 0x0a, 0x09,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4f, 0x0a, 0x0d, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69,
	0x65, 0x77, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x76, 0x69, 0x65, 0x77,
	0x49, 0x64, 0x12, 0x26, 0x0a, 0x0e, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x6c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x5b, 0x0a, 0x0b, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x70,
	0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x22, 0x3e, 0x0a, 0x0a, 0x46, 0x77, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x4f, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x73, 0x65,
	0x6e, 0x73, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xd3, 0x01, 0x0a, 0x0c, 0x4e, 0x6f, 0x64,
	0x65, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x40, 0x0a, 0x10, 0x43, 0x6f, 0x6e,
	0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x2e,
	0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75,
	0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x0b, 0x46,
	0x65, 0x74, 0x63, 0x68, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x16, 0x2e, 0x66, 0x61, 0x62,
	0x72, 0x69, 0x63, 0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x14, 0x2e, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0x00, 0x30, 0x01, 0x12, 0x40, 0x0a, 0x0f,
	0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x12, 0x2e, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x49, 0x44, 0x1a, 0x15, 0x2e, 0x66, 0x61, 0x62, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x2f, 0x61, 0x70, 0x70, 0x3b, 0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_messages_proto_rawDescOnce sync.Once
	file_node_messages_proto_rawDescData = file_node_messages_proto_rawDesc
)

func file_node_messages_proto_rawDescGZIP() []byte {
	file_node_messages_proto_rawDescOnce.Do(func() {
		file_node_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_messages_proto_rawDescData)
	})
	return file_node_messages_proto_rawDescData
}

var file_node_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_node_messages_proto_goTypes = []interface{}{
	(*ContentChunk)(nil),  // 0: fabrico.ContentChunk
	(*ContentID)(nil),     // 1: fabrico.ContentID
	(*BlockPosition)(nil), // 2: fabrico.BlockPosition
	(*BlockRecord)(nil),   // 3: fabrico.BlockRecord
	(*FwdMessage)(nil),    // 4: fabrico.FwdMessage
	(*Consensus)(nil),     // 5: fabrico.Consensus
	(*anypb.Any)(nil),     // 6: google.protobuf.Any
	(*emptypb.Empty)(nil), // 7: google.protobuf.Empty
}
var file_node_messages_proto_depIdxs = []int32{
	6, // 0: fabrico.Consensus.message:type_name -> google.protobuf.Any
	5, // 1: fabrico.NodeExchange.ConsensusMessage:input_type -> fabrico.Consensus
	2, // 2: fabrico.NodeExchange.FetchBlocks:input_type -> fabrico.BlockPosition
	1, // 3: fabrico.NodeExchange.DownloadContent:input_type -> fabrico.ContentID
	7, // 4: fabrico.NodeExchange.ConsensusMessage:output_type -> google.protobuf.Empty
	3, // 5: fabrico.NodeExchange.FetchBlocks:output_type -> fabrico.BlockRecord
	0, // 6: fabrico.NodeExchange.DownloadContent:output_type -> fabrico.ContentChunk
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_node_messages_proto_init() }
func file_node_messages_proto_init() {
	if File_node_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContentChunk); i {
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
		file_node_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContentID); i {
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
		file_node_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockPosition); i {
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
		file_node_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockRecord); i {
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
		file_node_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FwdMessage); i {
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
		file_node_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Consensus); i {
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
			RawDescriptor: file_node_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_node_messages_proto_goTypes,
		DependencyIndexes: file_node_messages_proto_depIdxs,
		MessageInfos:      file_node_messages_proto_msgTypes,
	}.Build()
	File_node_messages_proto = out.File
	file_node_messages_proto_rawDesc = nil
	file_node_messages_proto_goTypes = nil
	file_node_messages_proto_depIdxs = nil
}
