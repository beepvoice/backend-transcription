// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bite.proto

package main

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Bite struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Start                uint64   `protobuf:"varint,2,opt,name=start,proto3" json:"start,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Client               *Client  `protobuf:"bytes,4,opt,name=client,proto3" json:"client,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bite) Reset()         { *m = Bite{} }
func (m *Bite) String() string { return proto.CompactTextString(m) }
func (*Bite) ProtoMessage()    {}
func (*Bite) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1ec993646b17549, []int{0}
}

func (m *Bite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bite.Unmarshal(m, b)
}
func (m *Bite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bite.Marshal(b, m, deterministic)
}
func (m *Bite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bite.Merge(m, src)
}
func (m *Bite) XXX_Size() int {
	return xxx_messageInfo_Bite.Size(m)
}
func (m *Bite) XXX_DiscardUnknown() {
	xxx_messageInfo_Bite.DiscardUnknown(m)
}

var xxx_messageInfo_Bite proto.InternalMessageInfo

func (m *Bite) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Bite) GetStart() uint64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *Bite) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Bite) GetClient() *Client {
	if m != nil {
		return m.Client
	}
	return nil
}

func init() {
	proto.RegisterType((*Bite)(nil), "main.Bite")
}

func init() { proto.RegisterFile("bite.proto", fileDescriptor_e1ec993646b17549) }

var fileDescriptor_e1ec993646b17549 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xca, 0x2c, 0x49,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcc, 0xcc, 0x93, 0xe2, 0x49, 0xce,
	0xc9, 0x4c, 0xcd, 0x2b, 0x81, 0x88, 0x29, 0x65, 0x70, 0xb1, 0x38, 0x65, 0x96, 0xa4, 0x0a, 0x09,
	0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0x98, 0x42, 0x22,
	0x5c, 0xac, 0xc5, 0x25, 0x89, 0x45, 0x25, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x10, 0x8e,
	0x90, 0x10, 0x17, 0x4b, 0x4a, 0x62, 0x49, 0xa2, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x98,
	0x2d, 0xa4, 0xc2, 0xc5, 0x06, 0x31, 0x53, 0x82, 0x45, 0x81, 0x51, 0x83, 0xdb, 0x88, 0x47, 0x0f,
	0x64, 0x91, 0x9e, 0x33, 0x58, 0x2c, 0x08, 0x2a, 0x97, 0xc4, 0x06, 0xb6, 0xd0, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0x0a, 0x85, 0x92, 0x1b, 0x92, 0x00, 0x00, 0x00,
}