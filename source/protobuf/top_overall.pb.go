// Code generated by protoc-gen-go. DO NOT EDIT.
// source: top_overall.proto

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

type Entry struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Total                float32  `protobuf:"fixed32,2,opt,name=Total,proto3" json:"Total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_aeeca2cb3885364f, []int{0}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Entry) GetTotal() float32 {
	if m != nil {
		return m.Total
	}
	return 0
}

type TopOverallData struct {
	ID                   string             `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Year                 string             `protobuf:"bytes,2,opt,name=Year,proto3" json:"Year,omitempty"`
	Bucket               string             `protobuf:"bytes,3,opt,name=Bucket,proto3" json:"Bucket,omitempty"`
	Category             string             `protobuf:"bytes,4,opt,name=Category,proto3" json:"Category,omitempty"`
	Party                string             `protobuf:"bytes,5,opt,name=Party,proto3" json:"Party,omitempty"`
	Amts                 map[string]float32 `protobuf:"bytes,6,rep,name=Amts,proto3" json:"Amts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	Threshold            []*Entry           `protobuf:"bytes,7,rep,name=Threshold,proto3" json:"Threshold,omitempty"`
	SizeLimit            int32              `protobuf:"varint,8,opt,name=SizeLimit,proto3" json:"SizeLimit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TopOverallData) Reset()         { *m = TopOverallData{} }
func (m *TopOverallData) String() string { return proto.CompactTextString(m) }
func (*TopOverallData) ProtoMessage()    {}
func (*TopOverallData) Descriptor() ([]byte, []int) {
	return fileDescriptor_aeeca2cb3885364f, []int{1}
}

func (m *TopOverallData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopOverallData.Unmarshal(m, b)
}
func (m *TopOverallData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopOverallData.Marshal(b, m, deterministic)
}
func (m *TopOverallData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopOverallData.Merge(m, src)
}
func (m *TopOverallData) XXX_Size() int {
	return xxx_messageInfo_TopOverallData.Size(m)
}
func (m *TopOverallData) XXX_DiscardUnknown() {
	xxx_messageInfo_TopOverallData.DiscardUnknown(m)
}

var xxx_messageInfo_TopOverallData proto.InternalMessageInfo

func (m *TopOverallData) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *TopOverallData) GetYear() string {
	if m != nil {
		return m.Year
	}
	return ""
}

func (m *TopOverallData) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *TopOverallData) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *TopOverallData) GetParty() string {
	if m != nil {
		return m.Party
	}
	return ""
}

func (m *TopOverallData) GetAmts() map[string]float32 {
	if m != nil {
		return m.Amts
	}
	return nil
}

func (m *TopOverallData) GetThreshold() []*Entry {
	if m != nil {
		return m.Threshold
	}
	return nil
}

func (m *TopOverallData) GetSizeLimit() int32 {
	if m != nil {
		return m.SizeLimit
	}
	return 0
}

func init() {
	proto.RegisterType((*Entry)(nil), "protobuf.Entry")
	proto.RegisterType((*TopOverallData)(nil), "protobuf.TopOverallData")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.TopOverallData.AmtsEntry")
}

func init() { proto.RegisterFile("top_overall.proto", fileDescriptor_aeeca2cb3885364f) }

var fileDescriptor_aeeca2cb3885364f = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0xc9, 0xe6, 0x8f, 0xc9, 0x08, 0x55, 0x17, 0x91, 0xa5, 0x78, 0x08, 0x39, 0xe5, 0xd2,
	0x1c, 0x14, 0x54, 0xbc, 0xa9, 0xf1, 0x50, 0x10, 0x94, 0x98, 0x8b, 0x27, 0xd9, 0xea, 0x6a, 0x43,
	0xb7, 0x6e, 0xd8, 0x4e, 0x0a, 0xf1, 0xd9, 0x7c, 0x38, 0xc9, 0xa4, 0x6d, 0x10, 0x4f, 0x3b, 0xbf,
	0x99, 0xef, 0x9b, 0x9d, 0x19, 0x38, 0x42, 0x53, 0xbf, 0x9a, 0xb5, 0xb2, 0x52, 0xeb, 0xac, 0xb6,
	0x06, 0x0d, 0x0f, 0xe9, 0x99, 0x35, 0x1f, 0xc9, 0x04, 0xfc, 0xfb, 0x2f, 0xb4, 0x2d, 0x1f, 0x01,
	0x9b, 0xe6, 0xc2, 0x89, 0x9d, 0x34, 0x2a, 0xd8, 0x34, 0xe7, 0xc7, 0xe0, 0x97, 0x06, 0xa5, 0x16,
	0x2c, 0x76, 0x52, 0x56, 0xf4, 0x90, 0xfc, 0x30, 0x18, 0x95, 0xa6, 0x7e, 0xec, 0xbb, 0xe5, 0x12,
	0xe5, 0x3f, 0x23, 0x07, 0xef, 0x45, 0x49, 0x4b, 0xbe, 0xa8, 0xa0, 0x98, 0x9f, 0x40, 0x70, 0xdb,
	0xbc, 0x2d, 0x14, 0x0a, 0x97, 0xb2, 0x1b, 0xe2, 0x63, 0x08, 0xef, 0x24, 0xaa, 0x4f, 0x63, 0x5b,
	0xe1, 0x51, 0x65, 0xc7, 0xdd, 0x00, 0x4f, 0xd2, 0x62, 0x2b, 0x7c, 0x2a, 0xf4, 0xc0, 0x2f, 0xc0,
	0xbb, 0x59, 0xe2, 0x4a, 0x04, 0xb1, 0x9b, 0xee, 0x9f, 0x25, 0xd9, 0x76, 0x91, 0xec, 0xef, 0x54,
	0x59, 0x27, 0xa2, 0xc5, 0x0a, 0xd2, 0xf3, 0x09, 0x44, 0xe5, 0xdc, 0xaa, 0xd5, 0xdc, 0xe8, 0x77,
	0xb1, 0x47, 0xe6, 0x83, 0xc1, 0xdc, 0x2b, 0x07, 0x05, 0x3f, 0x85, 0xe8, 0xb9, 0xfa, 0x56, 0x0f,
	0xd5, 0xb2, 0x42, 0x11, 0xc6, 0x4e, 0xea, 0x17, 0x43, 0x62, 0x7c, 0x09, 0xd1, 0xae, 0x3f, 0x3f,
	0x04, 0x77, 0xa1, 0xda, 0xcd, 0x01, 0xba, 0xb0, 0x9b, 0x7c, 0x2d, 0x75, 0xa3, 0xb6, 0xa7, 0x23,
	0xb8, 0x66, 0x57, 0xce, 0x2c, 0xa0, 0x1f, 0xcf, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xd4, 0xe7,
	0xf2, 0x50, 0x93, 0x01, 0x00, 0x00,
}
