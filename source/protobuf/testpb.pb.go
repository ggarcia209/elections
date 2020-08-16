// Code generated by protoc-gen-go. DO NOT EDIT.
// source: testpb.proto

package protobuf

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

type Outer struct {
	Name                 string         `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Place                string         `protobuf:"bytes,2,opt,name=Place,proto3" json:"Place,omitempty"`
	Scores               []*Outer_Inner `protobuf:"bytes,3,rep,name=Scores,proto3" json:"Scores,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Outer) Reset()         { *m = Outer{} }
func (m *Outer) String() string { return proto.CompactTextString(m) }
func (*Outer) ProtoMessage()    {}
func (*Outer) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b98c0ed33edeb52, []int{0}
}

func (m *Outer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Outer.Unmarshal(m, b)
}
func (m *Outer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Outer.Marshal(b, m, deterministic)
}
func (m *Outer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Outer.Merge(m, src)
}
func (m *Outer) XXX_Size() int {
	return xxx_messageInfo_Outer.Size(m)
}
func (m *Outer) XXX_DiscardUnknown() {
	xxx_messageInfo_Outer.DiscardUnknown(m)
}

var xxx_messageInfo_Outer proto.InternalMessageInfo

func (m *Outer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Outer) GetPlace() string {
	if m != nil {
		return m.Place
	}
	return ""
}

func (m *Outer) GetScores() []*Outer_Inner {
	if m != nil {
		return m.Scores
	}
	return nil
}

type Outer_Inner struct {
	Score                int32    `protobuf:"varint,1,opt,name=Score,proto3" json:"Score,omitempty"`
	Grade                string   `protobuf:"bytes,2,opt,name=Grade,proto3" json:"Grade,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Outer_Inner) Reset()         { *m = Outer_Inner{} }
func (m *Outer_Inner) String() string { return proto.CompactTextString(m) }
func (*Outer_Inner) ProtoMessage()    {}
func (*Outer_Inner) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b98c0ed33edeb52, []int{0, 0}
}

func (m *Outer_Inner) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Outer_Inner.Unmarshal(m, b)
}
func (m *Outer_Inner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Outer_Inner.Marshal(b, m, deterministic)
}
func (m *Outer_Inner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Outer_Inner.Merge(m, src)
}
func (m *Outer_Inner) XXX_Size() int {
	return xxx_messageInfo_Outer_Inner.Size(m)
}
func (m *Outer_Inner) XXX_DiscardUnknown() {
	xxx_messageInfo_Outer_Inner.DiscardUnknown(m)
}

var xxx_messageInfo_Outer_Inner proto.InternalMessageInfo

func (m *Outer_Inner) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *Outer_Inner) GetGrade() string {
	if m != nil {
		return m.Grade
	}
	return ""
}

func init() {
	proto.RegisterType((*Outer)(nil), "protobuf.Outer")
	proto.RegisterType((*Outer_Inner)(nil), "protobuf.Outer.Inner")
}

func init() { proto.RegisterFile("testpb.proto", fileDescriptor_1b98c0ed33edeb52) }

var fileDescriptor_1b98c0ed33edeb52 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x49, 0x2d, 0x2e,
	0x29, 0x48, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a,
	0x53, 0x19, 0xb9, 0x58, 0xfd, 0x4b, 0x4b, 0x52, 0x8b, 0x84, 0x84, 0xb8, 0x58, 0xfc, 0x12, 0x73,
	0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x11, 0x2e, 0xd6, 0x80, 0x9c,
	0xc4, 0xe4, 0x54, 0x09, 0x26, 0xb0, 0x20, 0x84, 0x23, 0xa4, 0xcb, 0xc5, 0x16, 0x9c, 0x9c, 0x5f,
	0x94, 0x5a, 0x2c, 0xc1, 0xac, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xaa, 0x07, 0x33, 0x4e, 0x0f, 0x6c,
	0x94, 0x9e, 0x67, 0x5e, 0x5e, 0x6a, 0x51, 0x10, 0x54, 0x91, 0x94, 0x31, 0x17, 0x2b, 0x58, 0x00,
	0x64, 0x1a, 0x58, 0x08, 0x6c, 0x05, 0x6b, 0x10, 0x84, 0x03, 0x12, 0x75, 0x2f, 0x4a, 0x4c, 0x81,
	0xdb, 0x01, 0xe6, 0x24, 0xb1, 0x81, 0x8d, 0x34, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x28, 0xf0,
	0x37, 0xad, 0xb8, 0x00, 0x00, 0x00,
}