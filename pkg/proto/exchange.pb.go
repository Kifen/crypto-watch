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
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Exchange             string   `protobuf:"bytes,2,opt,name=Exchange,proto3" json:"Exchange,omitempty"`
	Req                  *Req     `protobuf:"bytes,3,opt,name=req,proto3" json:"req,omitempty"`
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

func (m *ExchangeReq) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ExchangeReq) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *ExchangeReq) GetReq() *Req {
	if m != nil {
		return m.Req
	}
	return nil
}

type ExchangeRes struct {
	Req                  *ExchangeReq `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	Price                float32      `protobuf:"fixed32,2,opt,name=price,proto3" json:"price,omitempty"`
	Message              string       `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
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

func (m *ExchangeRes) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Req struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Action               int32    `protobuf:"varint,2,opt,name=action,proto3" json:"action,omitempty"`
	Price                float32  `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_a678a54b61e3becc, []int{2}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Req) GetAction() int32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *Req) GetPrice() float32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func init() {
	proto.RegisterType((*ExchangeReq)(nil), "proto.ExchangeReq")
	proto.RegisterType((*ExchangeRes)(nil), "proto.ExchangeRes")
	proto.RegisterType((*Req)(nil), "proto.Req")
}

func init() {
	proto.RegisterFile("pkg/proto/exchange.proto", fileDescriptor_a678a54b61e3becc)
}

var fileDescriptor_a678a54b61e3becc = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8e, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xd9, 0x84, 0x54, 0x33, 0x11, 0x0f, 0x83, 0xc8, 0x52, 0x3c, 0x84, 0xe0, 0x21, 0xa7,
	0x56, 0xea, 0xcd, 0xab, 0x78, 0x10, 0x2f, 0x32, 0x97, 0x9e, 0xd3, 0xed, 0x90, 0x06, 0x6d, 0x37,
	0xc9, 0xae, 0x60, 0xff, 0xbd, 0x64, 0xb2, 0xd5, 0x80, 0x3d, 0x2d, 0xdf, 0x9b, 0xb7, 0xef, 0x3d,
	0xd0, 0xed, 0x47, 0xbd, 0x6c, 0x7b, 0xeb, 0xed, 0x92, 0xbf, 0xcd, 0xae, 0x3a, 0xd4, 0xbc, 0x10,
	0xc4, 0x44, 0x9e, 0x62, 0x0d, 0xd9, 0x4b, 0x38, 0x10, 0x77, 0x78, 0x0d, 0x51, 0xb3, 0xd5, 0x2a,
	0x57, 0x65, 0x42, 0x51, 0xb3, 0xc5, 0x39, 0x5c, 0x9e, 0xce, 0x3a, 0xca, 0x55, 0x99, 0xd2, 0x2f,
	0xe3, 0x1d, 0xc4, 0x3d, 0x77, 0x3a, 0xce, 0x55, 0x99, 0xad, 0x60, 0x8c, 0x5d, 0x10, 0x77, 0x34,
	0xc8, 0x85, 0x99, 0x06, 0x3b, 0xbc, 0x1f, 0xcd, 0x4a, 0xcc, 0x18, 0xcc, 0x93, 0x66, 0xf9, 0x84,
	0x37, 0x90, 0xb4, 0x7d, 0x63, 0xc6, 0xae, 0x88, 0x46, 0x40, 0x0d, 0x17, 0x7b, 0x76, 0xae, 0xaa,
	0x59, 0xca, 0x52, 0x3a, 0x61, 0xf1, 0x06, 0xf1, 0xb0, 0xfa, 0x16, 0x66, 0xee, 0xb8, 0xdf, 0xd8,
	0x4f, 0xc9, 0x4f, 0x29, 0xd0, 0xa0, 0x57, 0xc6, 0x37, 0xf6, 0x20, 0x79, 0x09, 0x05, 0xfa, 0xab,
	0x89, 0x27, 0x35, 0xab, 0x57, 0xc8, 0x9e, 0xfb, 0x63, 0xeb, 0xed, 0xba, 0xf2, 0x66, 0x87, 0x4f,
	0x70, 0x45, 0xdc, 0x7d, 0xb1, 0xf3, 0xef, 0xb2, 0xe2, 0xcc, 0xe8, 0xf9, 0x7f, 0xcd, 0x95, 0xea,
	0x41, 0x6d, 0x66, 0x22, 0x3f, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0xb8, 0x99, 0x16, 0x7c, 0x7f,
	0x01, 0x00, 0x00,
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
