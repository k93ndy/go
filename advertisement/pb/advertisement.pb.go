// Code generated by protoc-gen-go. DO NOT EDIT.
// source: advertisement.proto

package advertisement

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_advertisement_7ba1a80bec8b82f0, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type AdContent struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Image                []byte   `protobuf:"bytes,4,opt,name=image,proto3" json:"image,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdContent) Reset()         { *m = AdContent{} }
func (m *AdContent) String() string { return proto.CompactTextString(m) }
func (*AdContent) ProtoMessage()    {}
func (*AdContent) Descriptor() ([]byte, []int) {
	return fileDescriptor_advertisement_7ba1a80bec8b82f0, []int{1}
}
func (m *AdContent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdContent.Unmarshal(m, b)
}
func (m *AdContent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdContent.Marshal(b, m, deterministic)
}
func (dst *AdContent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdContent.Merge(dst, src)
}
func (m *AdContent) XXX_Size() int {
	return xxx_messageInfo_AdContent.Size(m)
}
func (m *AdContent) XXX_DiscardUnknown() {
	xxx_messageInfo_AdContent.DiscardUnknown(m)
}

var xxx_messageInfo_AdContent proto.InternalMessageInfo

func (m *AdContent) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AdContent) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *AdContent) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *AdContent) GetImage() []byte {
	if m != nil {
		return m.Image
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "advertisement.Empty")
	proto.RegisterType((*AdContent)(nil), "advertisement.AdContent")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AdvertisementClient is the client API for Advertisement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdvertisementClient interface {
	// request an advertisement info
	GetRandomAdvertisement(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AdContent, error)
}

type advertisementClient struct {
	cc *grpc.ClientConn
}

func NewAdvertisementClient(cc *grpc.ClientConn) AdvertisementClient {
	return &advertisementClient{cc}
}

func (c *advertisementClient) GetRandomAdvertisement(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AdContent, error) {
	out := new(AdContent)
	err := c.cc.Invoke(ctx, "/advertisement.Advertisement/getRandomAdvertisement", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdvertisementServer is the server API for Advertisement service.
type AdvertisementServer interface {
	// request an advertisement info
	GetRandomAdvertisement(context.Context, *Empty) (*AdContent, error)
}

func RegisterAdvertisementServer(s *grpc.Server, srv AdvertisementServer) {
	s.RegisterService(&_Advertisement_serviceDesc, srv)
}

func _Advertisement_GetRandomAdvertisement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServer).GetRandomAdvertisement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/advertisement.Advertisement/GetRandomAdvertisement",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServer).GetRandomAdvertisement(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Advertisement_serviceDesc = grpc.ServiceDesc{
	ServiceName: "advertisement.Advertisement",
	HandlerType: (*AdvertisementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getRandomAdvertisement",
			Handler:    _Advertisement_GetRandomAdvertisement_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "advertisement.proto",
}

func init() { proto.RegisterFile("advertisement.proto", fileDescriptor_advertisement_7ba1a80bec8b82f0) }

var fileDescriptor_advertisement_7ba1a80bec8b82f0 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0x4c, 0x29, 0x4b,
	0x2d, 0x2a, 0xc9, 0x2c, 0x4e, 0xcd, 0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x45, 0x11, 0x54, 0x62, 0xe7, 0x62, 0x75, 0xcd, 0x2d, 0x28, 0xa9, 0x54, 0xca, 0xe4, 0xe2,
	0x74, 0x4c, 0x71, 0xce, 0xcf, 0x2b, 0x49, 0xcd, 0x2b, 0x11, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc,
	0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85, 0x14, 0xb8, 0xb8, 0x53, 0x52,
	0x8b, 0x93, 0x8b, 0x32, 0x0b, 0x4a, 0x32, 0xf3, 0xf3, 0x24, 0x98, 0xc0, 0x52, 0xc8, 0x42, 0x42,
	0x02, 0x5c, 0xcc, 0xa5, 0x45, 0x39, 0x12, 0xcc, 0x60, 0x19, 0x10, 0x53, 0x48, 0x84, 0x8b, 0x35,
	0x33, 0x37, 0x31, 0x3d, 0x55, 0x82, 0x45, 0x81, 0x51, 0x83, 0x27, 0x08, 0xc2, 0x31, 0x8a, 0xe6,
	0xe2, 0x75, 0x44, 0x76, 0x84, 0x90, 0x17, 0x97, 0x58, 0x7a, 0x6a, 0x49, 0x50, 0x62, 0x5e, 0x4a,
	0x7e, 0x2e, 0xaa, 0x8c, 0x88, 0x1e, 0xaa, 0x1f, 0xc0, 0x6e, 0x95, 0x92, 0x40, 0x13, 0x85, 0x3b,
	0x5c, 0x89, 0x21, 0x89, 0x0d, 0xec, 0x4d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x22,
	0xf6, 0x06, 0xfd, 0x00, 0x00, 0x00,
}
