// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: node_messages.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeExchangeClient is the client API for NodeExchange service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeExchangeClient interface {
	ConsensusMessage(ctx context.Context, opts ...grpc.CallOption) (NodeExchange_ConsensusMessageClient, error)
	FetchBlocks(ctx context.Context, in *BlockPosition, opts ...grpc.CallOption) (NodeExchange_FetchBlocksClient, error)
	DownloadContent(ctx context.Context, in *ContentID, opts ...grpc.CallOption) (NodeExchange_DownloadContentClient, error)
}

type nodeExchangeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeExchangeClient(cc grpc.ClientConnInterface) NodeExchangeClient {
	return &nodeExchangeClient{cc}
}

func (c *nodeExchangeClient) ConsensusMessage(ctx context.Context, opts ...grpc.CallOption) (NodeExchange_ConsensusMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeExchange_ServiceDesc.Streams[0], "/fabrico.NodeExchange/ConsensusMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeExchangeConsensusMessageClient{stream}
	return x, nil
}

type NodeExchange_ConsensusMessageClient interface {
	Send(*Consensus) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type nodeExchangeConsensusMessageClient struct {
	grpc.ClientStream
}

func (x *nodeExchangeConsensusMessageClient) Send(m *Consensus) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeExchangeConsensusMessageClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeExchangeClient) FetchBlocks(ctx context.Context, in *BlockPosition, opts ...grpc.CallOption) (NodeExchange_FetchBlocksClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeExchange_ServiceDesc.Streams[1], "/fabrico.NodeExchange/FetchBlocks", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeExchangeFetchBlocksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NodeExchange_FetchBlocksClient interface {
	Recv() (*BlockRecord, error)
	grpc.ClientStream
}

type nodeExchangeFetchBlocksClient struct {
	grpc.ClientStream
}

func (x *nodeExchangeFetchBlocksClient) Recv() (*BlockRecord, error) {
	m := new(BlockRecord)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeExchangeClient) DownloadContent(ctx context.Context, in *ContentID, opts ...grpc.CallOption) (NodeExchange_DownloadContentClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeExchange_ServiceDesc.Streams[2], "/fabrico.NodeExchange/DownloadContent", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeExchangeDownloadContentClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NodeExchange_DownloadContentClient interface {
	Recv() (*ContentChunk, error)
	grpc.ClientStream
}

type nodeExchangeDownloadContentClient struct {
	grpc.ClientStream
}

func (x *nodeExchangeDownloadContentClient) Recv() (*ContentChunk, error) {
	m := new(ContentChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeExchangeServer is the server API for NodeExchange service.
// All implementations must embed UnimplementedNodeExchangeServer
// for forward compatibility
type NodeExchangeServer interface {
	ConsensusMessage(NodeExchange_ConsensusMessageServer) error
	FetchBlocks(*BlockPosition, NodeExchange_FetchBlocksServer) error
	DownloadContent(*ContentID, NodeExchange_DownloadContentServer) error
	mustEmbedUnimplementedNodeExchangeServer()
}

// UnimplementedNodeExchangeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeExchangeServer struct {
}

func (UnimplementedNodeExchangeServer) ConsensusMessage(NodeExchange_ConsensusMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method ConsensusMessage not implemented")
}
func (UnimplementedNodeExchangeServer) FetchBlocks(*BlockPosition, NodeExchange_FetchBlocksServer) error {
	return status.Errorf(codes.Unimplemented, "method FetchBlocks not implemented")
}
func (UnimplementedNodeExchangeServer) DownloadContent(*ContentID, NodeExchange_DownloadContentServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadContent not implemented")
}
func (UnimplementedNodeExchangeServer) mustEmbedUnimplementedNodeExchangeServer() {}

// UnsafeNodeExchangeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeExchangeServer will
// result in compilation errors.
type UnsafeNodeExchangeServer interface {
	mustEmbedUnimplementedNodeExchangeServer()
}

func RegisterNodeExchangeServer(s grpc.ServiceRegistrar, srv NodeExchangeServer) {
	s.RegisterService(&NodeExchange_ServiceDesc, srv)
}

func _NodeExchange_ConsensusMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeExchangeServer).ConsensusMessage(&nodeExchangeConsensusMessageServer{stream})
}

type NodeExchange_ConsensusMessageServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*Consensus, error)
	grpc.ServerStream
}

type nodeExchangeConsensusMessageServer struct {
	grpc.ServerStream
}

func (x *nodeExchangeConsensusMessageServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeExchangeConsensusMessageServer) Recv() (*Consensus, error) {
	m := new(Consensus)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _NodeExchange_FetchBlocks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BlockPosition)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeExchangeServer).FetchBlocks(m, &nodeExchangeFetchBlocksServer{stream})
}

type NodeExchange_FetchBlocksServer interface {
	Send(*BlockRecord) error
	grpc.ServerStream
}

type nodeExchangeFetchBlocksServer struct {
	grpc.ServerStream
}

func (x *nodeExchangeFetchBlocksServer) Send(m *BlockRecord) error {
	return x.ServerStream.SendMsg(m)
}

func _NodeExchange_DownloadContent_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ContentID)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeExchangeServer).DownloadContent(m, &nodeExchangeDownloadContentServer{stream})
}

type NodeExchange_DownloadContentServer interface {
	Send(*ContentChunk) error
	grpc.ServerStream
}

type nodeExchangeDownloadContentServer struct {
	grpc.ServerStream
}

func (x *nodeExchangeDownloadContentServer) Send(m *ContentChunk) error {
	return x.ServerStream.SendMsg(m)
}

// NodeExchange_ServiceDesc is the grpc.ServiceDesc for NodeExchange service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeExchange_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fabrico.NodeExchange",
	HandlerType: (*NodeExchangeServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConsensusMessage",
			Handler:       _NodeExchange_ConsensusMessage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "FetchBlocks",
			Handler:       _NodeExchange_FetchBlocks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadContent",
			Handler:       _NodeExchange_DownloadContent_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "node_messages.proto",
}
