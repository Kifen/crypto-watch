// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/proto/exchange.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ExchangeReq struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Exchange             string   `protobuf:"bytes,2,opt,name=Exchange,proto3" json:"Exchange,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExchangeReq) Reset()         { *m = ExchangeReq{} }
func (m *ExchangeReq) String() string { return proto.CompactTextString(m) }
func (*ExchangeReq) ProtoMessage()    {}
func (*ExchangeReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a678a54b61e3becc, []int{0}
}

func (m *ExchangeReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExchangeReq.Unmarshal(m, b)
}
func (m *ExchangeReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExchangeReq.Marshal(b, m, deterministic)
}
func (m *ExchangeReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeReq.Merge(m, src)
}
func (m *ExchangeReq) XXX_Size() int {
	return xxx_messageInfo_ExchangeReq.Size(m)
}
func (m *ExchangeReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeReq.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeReq proto.InternalMessageInfo

func (m *ExchangeReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExchangeReq) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

type ExchangeRes struct {
	Req                  *ExchangeReq `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	Price                float32      `protobuf:"fixed32,2,opt,name=price,proto3" json:"price,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ExchangeRes) Reset()         { *m = ExchangeRes{} }
func (m *ExchangeRes) String() string { return proto.CompactTextString(m) }
func (*ExchangeRes) ProtoMessage()    {}
func (*ExchangeRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_a678a54b61e3becc, []int{1}
}

func (m *ExchangeRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExchangeRes.Unmarshal(m, b)
}
func (m *ExchangeRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExchangeRes.Marshal(b, m, deterministic)
}
func (m *ExchangeRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeRes.Merge(m, src)
}
func (m *ExchangeRes) XXX_Size() int {
	return xxx_messageInfo_ExchangeRes.Size(m)
}
func (m *ExchangeRes) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeRes.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeRes proto.InternalMessageInfo

func (m *ExchangeRes) GetReq() *ExchangeReq {
	if m != nil {
		return m.Req
	}
	return nil
}

func (m *ExchangeRes) GetPrice() float32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func init() {
	proto.RegisterType((*ExchangeReq)(nil), "proto.ExchangeReq")
	proto.RegisterType((*ExchangeRes)(nil), "proto.ExchangeRes")
}

func init() {
	proto.RegisterFile("pkg/proto/exchange.proto", fileDescriptor_a678a54b61e3becc)
}

var fileDescriptor_a678a54b61e3becc = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0xc8, 0x4e, 0xd7,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0xad, 0x48, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0xd5, 0x03,
	0x73, 0x85, 0x58, 0xc1, 0x94, 0x92, 0x25, 0x17, 0xb7, 0x2b, 0x54, 0x22, 0x28, 0xb5, 0x50, 0x88,
	0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x88, 0x29, 0x33, 0x45, 0x48,
	0x8a, 0x8b, 0x03, 0x26, 0x2d, 0xc1, 0x04, 0x16, 0x85, 0xf3, 0x95, 0x3c, 0x91, 0xb5, 0x16, 0x0b,
	0xa9, 0x70, 0x31, 0x17, 0xa5, 0x16, 0x82, 0xf5, 0x72, 0x1b, 0x09, 0x41, 0x6c, 0xd1, 0x43, 0x32,
	0x3b, 0x08, 0x24, 0x2d, 0x24, 0xc2, 0xc5, 0x5a, 0x50, 0x94, 0x99, 0x0c, 0x31, 0x8d, 0x29, 0x08,
	0xc2, 0x31, 0xf2, 0xe4, 0xe2, 0x76, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0x0f, 0x4f, 0x2c, 0x49, 0xce,
	0x10, 0xb2, 0xe2, 0xe2, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x09, 0x00, 0x49, 0x0b, 0x61,
	0x31, 0x4d, 0x0a, 0x53, 0xac, 0x58, 0x83, 0xd1, 0x80, 0x31, 0x89, 0x0d, 0x2c, 0x6c, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xaa, 0xc3, 0xd8, 0xfd, 0xfa, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CryptoWatchClient is the client API for CryptoWatch service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CryptoWatchClient interface {
	RequestPrice(ctx context.Context, opts ...grpc.CallOption) (CryptoWatch_RequestPriceClient, error)
}

type cryptoWatchClient struct {
	cc grpc.ClientConnInterface
}

func NewCryptoWatchClient(cc grpc.ClientConnInterface) CryptoWatchClient {
	return &cryptoWatchClient{cc}
}

func (c *cryptoWatchClient) RequestPrice(ctx context.Context, opts ...grpc.CallOption) (CryptoWatch_RequestPriceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_CryptoWatch_serviceDesc.Streams[0], "/proto.CryptoWatch/RequestPrice", opts...)
	if err != nil {
		return nil, err
	}
	x := &cryptoWatchRequestPriceClient{stream}
	return x, nil
}

type CryptoWatch_RequestPriceClient interface {
	Send(*ExchangeReq) error
	Recv() (*ExchangeRes, error)
	grpc.ClientStream
}

type cryptoWatchRequestPriceClient struct {
	grpc.ClientStream
}

func (x *cryptoWatchRequestPriceClient) Send(m *ExchangeReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cryptoWatchRequestPriceClient) Recv() (*ExchangeRes, error) {
	m := new(ExchangeRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CryptoWatchServer is the server API for CryptoWatch service.
type CryptoWatchServer interface {
	RequestPrice(CryptoWatch_RequestPriceServer) error
}

// UnimplementedCryptoWatchServer can be embedded to have forward compatible implementations.
type UnimplementedCryptoWatchServer struct {
}

func (*UnimplementedCryptoWatchServer) RequestPrice(srv CryptoWatch_RequestPriceServer) error {
	return status.Errorf(codes.Unimplemented, "method RequestPrice not implemented")
}

func RegisterCryptoWatchServer(s *grpc.Server, srv CryptoWatchServer) {
	s.RegisterService(&_CryptoWatch_serviceDesc, srv)
}

func _CryptoWatch_RequestPrice_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CryptoWatchServer).RequestPrice(&cryptoWatchRequestPriceServer{stream})
}

type CryptoWatch_RequestPriceServer interface {
	Send(*ExchangeRes) error
	Recv() (*ExchangeReq, error)
	grpc.ServerStream
}

type cryptoWatchRequestPriceServer struct {
	grpc.ServerStream
}

func (x *cryptoWatchRequestPriceServer) Send(m *ExchangeRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cryptoWatchRequestPriceServer) Recv() (*ExchangeReq, error) {
	m := new(ExchangeReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _CryptoWatch_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CryptoWatch",
	HandlerType: (*CryptoWatchServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RequestPrice",
			Handler:       _CryptoWatch_RequestPrice_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/proto/exchange.proto",
}
