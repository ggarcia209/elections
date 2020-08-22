// Code generated by protoc-gen-go. DO NOT EDIT.
// source: index.proto

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

type SearchResult struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	City                 string   `protobuf:"bytes,3,opt,name=City,proto3" json:"City,omitempty"`
	State                string   `protobuf:"bytes,4,opt,name=State,proto3" json:"State,omitempty"`
	Bucket               string   `protobuf:"bytes,5,opt,name=Bucket,proto3" json:"Bucket,omitempty"`
	Years                []string `protobuf:"bytes,6,rep,name=Years,proto3" json:"Years,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchResult) Reset()         { *m = SearchResult{} }
func (m *SearchResult) String() string { return proto.CompactTextString(m) }
func (*SearchResult) ProtoMessage()    {}
func (*SearchResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{0}
}

func (m *SearchResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchResult.Unmarshal(m, b)
}
func (m *SearchResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchResult.Marshal(b, m, deterministic)
}
func (m *SearchResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchResult.Merge(m, src)
}
func (m *SearchResult) XXX_Size() int {
	return xxx_messageInfo_SearchResult.Size(m)
}
func (m *SearchResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchResult.DiscardUnknown(m)
}

var xxx_messageInfo_SearchResult proto.InternalMessageInfo

func (m *SearchResult) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *SearchResult) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SearchResult) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *SearchResult) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *SearchResult) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *SearchResult) GetYears() []string {
	if m != nil {
		return m.Years
	}
	return nil
}

type ResultList struct {
	IDs                  []string `protobuf:"bytes,1,rep,name=IDs,proto3" json:"IDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResultList) Reset()         { *m = ResultList{} }
func (m *ResultList) String() string { return proto.CompactTextString(m) }
func (*ResultList) ProtoMessage()    {}
func (*ResultList) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{1}
}

func (m *ResultList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResultList.Unmarshal(m, b)
}
func (m *ResultList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResultList.Marshal(b, m, deterministic)
}
func (m *ResultList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResultList.Merge(m, src)
}
func (m *ResultList) XXX_Size() int {
	return xxx_messageInfo_ResultList.Size(m)
}
func (m *ResultList) XXX_DiscardUnknown() {
	xxx_messageInfo_ResultList.DiscardUnknown(m)
}

var xxx_messageInfo_ResultList proto.InternalMessageInfo

func (m *ResultList) GetIDs() []string {
	if m != nil {
		return m.IDs
	}
	return nil
}

type LookupMap struct {
	Lookup               map[string]*SearchResult `protobuf:"bytes,1,rep,name=Lookup,proto3" json:"Lookup,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *LookupMap) Reset()         { *m = LookupMap{} }
func (m *LookupMap) String() string { return proto.CompactTextString(m) }
func (*LookupMap) ProtoMessage()    {}
func (*LookupMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{2}
}

func (m *LookupMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LookupMap.Unmarshal(m, b)
}
func (m *LookupMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LookupMap.Marshal(b, m, deterministic)
}
func (m *LookupMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LookupMap.Merge(m, src)
}
func (m *LookupMap) XXX_Size() int {
	return xxx_messageInfo_LookupMap.Size(m)
}
func (m *LookupMap) XXX_DiscardUnknown() {
	xxx_messageInfo_LookupMap.DiscardUnknown(m)
}

var xxx_messageInfo_LookupMap proto.InternalMessageInfo

func (m *LookupMap) GetLookup() map[string]*SearchResult {
	if m != nil {
		return m.Lookup
	}
	return nil
}

type IndexData struct {
	Size                 float32              `protobuf:"fixed32,1,opt,name=Size,proto3" json:"Size,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,2,opt,name=LastUpdated,proto3" json:"LastUpdated,omitempty"`
	Completed            map[string]bool      `protobuf:"bytes,3,rep,name=Completed,proto3" json:"Completed,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *IndexData) Reset()         { *m = IndexData{} }
func (m *IndexData) String() string { return proto.CompactTextString(m) }
func (*IndexData) ProtoMessage()    {}
func (*IndexData) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{3}
}

func (m *IndexData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexData.Unmarshal(m, b)
}
func (m *IndexData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexData.Marshal(b, m, deterministic)
}
func (m *IndexData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexData.Merge(m, src)
}
func (m *IndexData) XXX_Size() int {
	return xxx_messageInfo_IndexData.Size(m)
}
func (m *IndexData) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexData.DiscardUnknown(m)
}

var xxx_messageInfo_IndexData proto.InternalMessageInfo

func (m *IndexData) GetSize() float32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *IndexData) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func (m *IndexData) GetCompleted() map[string]bool {
	if m != nil {
		return m.Completed
	}
	return nil
}

type PartitionMap struct {
	Partitions           map[string]bool `protobuf:"bytes,1,rep,name=Partitions,proto3" json:"Partitions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PartitionMap) Reset()         { *m = PartitionMap{} }
func (m *PartitionMap) String() string { return proto.CompactTextString(m) }
func (*PartitionMap) ProtoMessage()    {}
func (*PartitionMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{4}
}

func (m *PartitionMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PartitionMap.Unmarshal(m, b)
}
func (m *PartitionMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PartitionMap.Marshal(b, m, deterministic)
}
func (m *PartitionMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartitionMap.Merge(m, src)
}
func (m *PartitionMap) XXX_Size() int {
	return xxx_messageInfo_PartitionMap.Size(m)
}
func (m *PartitionMap) XXX_DiscardUnknown() {
	xxx_messageInfo_PartitionMap.DiscardUnknown(m)
}

var xxx_messageInfo_PartitionMap proto.InternalMessageInfo

func (m *PartitionMap) GetPartitions() map[string]bool {
	if m != nil {
		return m.Partitions
	}
	return nil
}

func init() {
	proto.RegisterType((*SearchResult)(nil), "protobuf.SearchResult")
	proto.RegisterType((*ResultList)(nil), "protobuf.ResultList")
	proto.RegisterType((*LookupMap)(nil), "protobuf.LookupMap")
	proto.RegisterMapType((map[string]*SearchResult)(nil), "protobuf.LookupMap.LookupEntry")
	proto.RegisterType((*IndexData)(nil), "protobuf.IndexData")
	proto.RegisterMapType((map[string]bool)(nil), "protobuf.IndexData.CompletedEntry")
	proto.RegisterType((*PartitionMap)(nil), "protobuf.PartitionMap")
	proto.RegisterMapType((map[string]bool)(nil), "protobuf.PartitionMap.PartitionsEntry")
}

func init() { proto.RegisterFile("index.proto", fileDescriptor_f750e0f7889345b5) }

var fileDescriptor_f750e0f7889345b5 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x51, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0x65, 0x25, 0x5b, 0x58, 0x23, 0xe3, 0x96, 0xc5, 0x18, 0xa1, 0x43, 0x6d, 0x74, 0x28, 0x3d,
	0x14, 0x19, 0xdc, 0x43, 0x4b, 0x71, 0xa1, 0xd4, 0x6a, 0x41, 0xe0, 0x96, 0x56, 0x6e, 0x0e, 0x39,
	0xae, 0xed, 0x8d, 0x23, 0x2c, 0x79, 0x85, 0x76, 0x15, 0xe2, 0x7c, 0x42, 0xee, 0x21, 0x5f, 0x97,
	0x7f, 0x09, 0xbb, 0x2b, 0x4b, 0x4a, 0xc8, 0x25, 0x27, 0xbd, 0xf7, 0xf4, 0x66, 0x67, 0xde, 0x0c,
	0x38, 0xc9, 0x61, 0x4b, 0xaf, 0x83, 0xbc, 0x60, 0x82, 0xe1, 0x9e, 0xfa, 0xac, 0xcb, 0x0b, 0x6f,
	0xbc, 0x63, 0x6c, 0x97, 0xd2, 0xe9, 0x49, 0x98, 0x8a, 0x24, 0xa3, 0x5c, 0x90, 0x2c, 0xd7, 0x56,
	0xff, 0x16, 0x41, 0x7f, 0x45, 0x49, 0xb1, 0xb9, 0x8c, 0x29, 0x2f, 0x53, 0x81, 0x07, 0x60, 0x44,
	0xa1, 0x8b, 0x26, 0xe8, 0x83, 0x1d, 0x1b, 0x51, 0x88, 0x31, 0x74, 0xfe, 0x90, 0x8c, 0xba, 0x86,
	0x52, 0x14, 0x96, 0xda, 0x22, 0x11, 0x47, 0xd7, 0xd4, 0x9a, 0xc4, 0x78, 0x08, 0xdd, 0x95, 0x20,
	0x82, 0xba, 0x1d, 0x25, 0x6a, 0x82, 0x47, 0x60, 0xfd, 0x28, 0x37, 0x7b, 0x2a, 0xdc, 0xae, 0x92,
	0x2b, 0x26, 0xdd, 0xe7, 0x94, 0x14, 0xdc, 0xb5, 0x26, 0xa6, 0x74, 0x2b, 0xe2, 0xbf, 0x03, 0xd0,
	0x53, 0x2c, 0x13, 0x2e, 0xf0, 0x5b, 0x30, 0xa3, 0x90, 0xbb, 0x48, 0x39, 0x24, 0xf4, 0xef, 0x11,
	0xd8, 0x4b, 0xc6, 0xf6, 0x65, 0xfe, 0x9b, 0xe4, 0xf8, 0x33, 0x58, 0x9a, 0x28, 0x8b, 0x33, 0x1b,
	0x07, 0xa7, 0x94, 0x41, 0x6d, 0xaa, 0xd0, 0xcf, 0x83, 0x28, 0x8e, 0x71, 0x65, 0xf7, 0xfe, 0x81,
	0xd3, 0x92, 0x65, 0x9f, 0x3d, 0x3d, 0x56, 0x91, 0x25, 0xc4, 0x1f, 0xa1, 0x7b, 0x45, 0xd2, 0x52,
	0x87, 0x76, 0x66, 0xa3, 0xe6, 0xe1, 0xf6, 0xaa, 0x62, 0x6d, 0xfa, 0x6a, 0x7c, 0x41, 0xfe, 0x03,
	0x02, 0x3b, 0x92, 0x17, 0x08, 0x89, 0x20, 0x72, 0x3f, 0xab, 0xe4, 0x86, 0xaa, 0x27, 0x8d, 0x58,
	0x61, 0x3c, 0x07, 0x67, 0x49, 0xb8, 0x38, 0xcb, 0xb7, 0x44, 0xd0, 0x6d, 0xf5, 0xb2, 0x17, 0xe8,
	0xfb, 0x34, 0x0d, 0xfe, 0x9f, 0xee, 0x13, 0xb7, 0xed, 0xf8, 0x3b, 0xd8, 0x0b, 0x96, 0xe5, 0x29,
	0x95, 0xb5, 0xa6, 0x8a, 0xeb, 0x37, 0x45, 0x75, 0xe7, 0xa0, 0x36, 0xe9, 0xc4, 0x4d, 0x91, 0x37,
	0x87, 0xc1, 0xd3, 0x9f, 0x2f, 0xe4, 0x1e, 0xb6, 0x73, 0xf7, 0xda, 0xf9, 0xee, 0x10, 0xf4, 0xff,
	0x92, 0x42, 0x24, 0x22, 0x61, 0x07, 0xb9, 0xfc, 0x5f, 0x00, 0x35, 0xe7, 0xd5, 0x01, 0xde, 0x37,
	0x13, 0xb5, 0xbd, 0x0d, 0xe1, 0x7a, 0xaa, 0x56, 0xa5, 0xf7, 0x0d, 0xde, 0x3c, 0xfb, 0xfd, 0x9a,
	0xb9, 0xd6, 0x96, 0xea, 0xf8, 0xe9, 0x31, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x24, 0x47, 0xaa, 0xff,
	0x02, 0x00, 0x00,
}
