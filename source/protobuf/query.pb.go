// Code generated by protoc-gen-go. DO NOT EDIT.
// source: query.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Query struct {
	Text                 string               `protobuf:"bytes,1,opt,name=Text,proto3" json:"Text,omitempty"`
	UserID               string               `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	TimeStamp            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Query) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Query) GetTimeStamp() *timestamp.Timestamp {
	if m != nil {
		return m.TimeStamp
	}
	return nil
}

func init() {
	proto.RegisterType((*Query)(nil), "protobuf.Query")
}

func init() { proto.RegisterFile("query.proto", fileDescriptor_5c6ac9b241082464) }

var fileDescriptor_5c6ac9b241082464 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2c, 0x4d, 0x2d,
	0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x52, 0xf2,
	0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0xfa, 0x30, 0x01, 0xfd, 0x92, 0xcc, 0xdc, 0xd4, 0xe2, 0x92,
	0xc4, 0xdc, 0x02, 0x88, 0x52, 0xa5, 0x5c, 0x2e, 0xd6, 0x40, 0x90, 0x4e, 0x21, 0x21, 0x2e, 0x96,
	0x90, 0xd4, 0x8a, 0x12, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b, 0x48, 0x8c, 0x8b,
	0x2d, 0xb4, 0x38, 0xb5, 0xc8, 0xd3, 0x45, 0x82, 0x09, 0x2c, 0x0a, 0xe5, 0x09, 0x59, 0x70, 0x71,
	0x86, 0x64, 0xe6, 0xa6, 0x06, 0x83, 0xcc, 0x91, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x36, 0x92, 0xd2,
	0x83, 0xd8, 0xa4, 0x07, 0xb3, 0x49, 0x2f, 0x04, 0x66, 0x53, 0x10, 0x42, 0x71, 0x12, 0x1b, 0x58,
	0xda, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xef, 0x56, 0x62, 0x49, 0xaf, 0x00, 0x00, 0x00,
}
