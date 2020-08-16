// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmte_cont.proto

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

type CmteContribution struct {
	CmteID               string               `protobuf:"bytes,1,opt,name=CmteID,proto3" json:"CmteID,omitempty"`
	AmndtInd             string               `protobuf:"bytes,2,opt,name=AmndtInd,proto3" json:"AmndtInd,omitempty"`
	ReportType           string               `protobuf:"bytes,3,opt,name=ReportType,proto3" json:"ReportType,omitempty"`
	TxPGI                string               `protobuf:"bytes,4,opt,name=TxPGI,proto3" json:"TxPGI,omitempty"`
	ImgNum               string               `protobuf:"bytes,5,opt,name=imgNum,proto3" json:"imgNum,omitempty"`
	TxType               string               `protobuf:"bytes,6,opt,name=TxType,proto3" json:"TxType,omitempty"`
	EntityType           string               `protobuf:"bytes,7,opt,name=EntityType,proto3" json:"EntityType,omitempty"`
	Name                 string               `protobuf:"bytes,8,opt,name=Name,proto3" json:"Name,omitempty"`
	City                 string               `protobuf:"bytes,9,opt,name=City,proto3" json:"City,omitempty"`
	State                string               `protobuf:"bytes,10,opt,name=State,proto3" json:"State,omitempty"`
	Zip                  string               `protobuf:"bytes,11,opt,name=Zip,proto3" json:"Zip,omitempty"`
	Employer             string               `protobuf:"bytes,12,opt,name=Employer,proto3" json:"Employer,omitempty"`
	Occupation           string               `protobuf:"bytes,13,opt,name=Occupation,proto3" json:"Occupation,omitempty"`
	TxDate               *timestamp.Timestamp `protobuf:"bytes,14,opt,name=TxDate,proto3" json:"TxDate,omitempty"`
	TxAmt                float32              `protobuf:"fixed32,15,opt,name=TxAmt,proto3" json:"TxAmt,omitempty"`
	OtherID              string               `protobuf:"bytes,16,opt,name=OtherID,proto3" json:"OtherID,omitempty"`
	CandID               string               `protobuf:"bytes,17,opt,name=CandID,proto3" json:"CandID,omitempty"`
	TxID                 string               `protobuf:"bytes,18,opt,name=TxID,proto3" json:"TxID,omitempty"`
	FileNum              int32                `protobuf:"varint,19,opt,name=FileNum,proto3" json:"FileNum,omitempty"`
	MemoCode             string               `protobuf:"bytes,20,opt,name=MemoCode,proto3" json:"MemoCode,omitempty"`
	MemoText             string               `protobuf:"bytes,21,opt,name=MemoText,proto3" json:"MemoText,omitempty"`
	SubID                int32                `protobuf:"varint,22,opt,name=SubID,proto3" json:"SubID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CmteContribution) Reset()         { *m = CmteContribution{} }
func (m *CmteContribution) String() string { return proto.CompactTextString(m) }
func (*CmteContribution) ProtoMessage()    {}
func (*CmteContribution) Descriptor() ([]byte, []int) {
	return fileDescriptor_de83e7a30dc43abf, []int{0}
}

func (m *CmteContribution) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmteContribution.Unmarshal(m, b)
}
func (m *CmteContribution) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmteContribution.Marshal(b, m, deterministic)
}
func (m *CmteContribution) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmteContribution.Merge(m, src)
}
func (m *CmteContribution) XXX_Size() int {
	return xxx_messageInfo_CmteContribution.Size(m)
}
func (m *CmteContribution) XXX_DiscardUnknown() {
	xxx_messageInfo_CmteContribution.DiscardUnknown(m)
}

var xxx_messageInfo_CmteContribution proto.InternalMessageInfo

func (m *CmteContribution) GetCmteID() string {
	if m != nil {
		return m.CmteID
	}
	return ""
}

func (m *CmteContribution) GetAmndtInd() string {
	if m != nil {
		return m.AmndtInd
	}
	return ""
}

func (m *CmteContribution) GetReportType() string {
	if m != nil {
		return m.ReportType
	}
	return ""
}

func (m *CmteContribution) GetTxPGI() string {
	if m != nil {
		return m.TxPGI
	}
	return ""
}

func (m *CmteContribution) GetImgNum() string {
	if m != nil {
		return m.ImgNum
	}
	return ""
}

func (m *CmteContribution) GetTxType() string {
	if m != nil {
		return m.TxType
	}
	return ""
}

func (m *CmteContribution) GetEntityType() string {
	if m != nil {
		return m.EntityType
	}
	return ""
}

func (m *CmteContribution) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CmteContribution) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *CmteContribution) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *CmteContribution) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *CmteContribution) GetEmployer() string {
	if m != nil {
		return m.Employer
	}
	return ""
}

func (m *CmteContribution) GetOccupation() string {
	if m != nil {
		return m.Occupation
	}
	return ""
}

func (m *CmteContribution) GetTxDate() *timestamp.Timestamp {
	if m != nil {
		return m.TxDate
	}
	return nil
}

func (m *CmteContribution) GetTxAmt() float32 {
	if m != nil {
		return m.TxAmt
	}
	return 0
}

func (m *CmteContribution) GetOtherID() string {
	if m != nil {
		return m.OtherID
	}
	return ""
}

func (m *CmteContribution) GetCandID() string {
	if m != nil {
		return m.CandID
	}
	return ""
}

func (m *CmteContribution) GetTxID() string {
	if m != nil {
		return m.TxID
	}
	return ""
}

func (m *CmteContribution) GetFileNum() int32 {
	if m != nil {
		return m.FileNum
	}
	return 0
}

func (m *CmteContribution) GetMemoCode() string {
	if m != nil {
		return m.MemoCode
	}
	return ""
}

func (m *CmteContribution) GetMemoText() string {
	if m != nil {
		return m.MemoText
	}
	return ""
}

func (m *CmteContribution) GetSubID() int32 {
	if m != nil {
		return m.SubID
	}
	return 0
}

func init() {
	proto.RegisterType((*CmteContribution)(nil), "protobuf.CmteContribution")
}

func init() { proto.RegisterFile("cmte_cont.proto", fileDescriptor_de83e7a30dc43abf) }

var fileDescriptor_de83e7a30dc43abf = []byte{
	// 389 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe5, 0x36, 0x49, 0xd3, 0x2d, 0xd0, 0xb0, 0x94, 0x6a, 0x94, 0x03, 0x44, 0x9c, 0x7c,
	0x72, 0xa5, 0xf2, 0x04, 0x95, 0x5d, 0x90, 0x0f, 0xb4, 0xc8, 0xf8, 0xc4, 0x05, 0xf9, 0xcf, 0x12,
	0x56, 0xca, 0x7a, 0x2d, 0x33, 0x96, 0xec, 0xb7, 0xe4, 0x91, 0xd0, 0xcc, 0x64, 0x09, 0x27, 0xcf,
	0xef, 0xfb, 0xac, 0x9d, 0xf9, 0x66, 0xd4, 0x75, 0xe3, 0xd0, 0xfc, 0x68, 0x7c, 0x87, 0x49, 0x3f,
	0x78, 0xf4, 0x7a, 0xcd, 0x9f, 0x7a, 0xfc, 0xb9, 0x7d, 0xbf, 0xf7, 0x7e, 0x7f, 0x30, 0x77, 0x41,
	0xb8, 0x43, 0xeb, 0xcc, 0x6f, 0xac, 0x5c, 0x2f, 0xbf, 0x7e, 0xf8, 0xb3, 0x50, 0x9b, 0xd4, 0xa1,
	0x49, 0x7d, 0x87, 0x83, 0xad, 0x47, 0xb4, 0xbe, 0xd3, 0xb7, 0x6a, 0x45, 0x5a, 0x9e, 0x41, 0xb4,
	0x8b, 0xe2, 0xcb, 0xe2, 0x48, 0x7a, 0xab, 0xd6, 0x0f, 0xae, 0x6b, 0x31, 0xef, 0x5a, 0x38, 0x63,
	0xe7, 0x1f, 0xeb, 0x77, 0x4a, 0x15, 0xa6, 0xf7, 0x03, 0x96, 0x73, 0x6f, 0xe0, 0x9c, 0xdd, 0xff,
	0x14, 0x7d, 0xa3, 0x96, 0xe5, 0xf4, 0xf5, 0x73, 0x0e, 0x0b, 0xb6, 0x04, 0xa8, 0x93, 0x75, 0xfb,
	0xa7, 0xd1, 0xc1, 0x52, 0x3a, 0x09, 0x91, 0x5e, 0x4e, 0xfc, 0xd2, 0x4a, 0x74, 0x21, 0xea, 0xf2,
	0xd8, 0xa1, 0xc5, 0x99, 0xbd, 0x0b, 0xe9, 0x72, 0x52, 0xb4, 0x56, 0x8b, 0xa7, 0xca, 0x19, 0x58,
	0xb3, 0xc3, 0x35, 0x69, 0xa9, 0xc5, 0x19, 0x2e, 0x45, 0xa3, 0x9a, 0xa6, 0xf9, 0x86, 0x15, 0x1a,
	0x50, 0x32, 0x0d, 0x83, 0xde, 0xa8, 0xf3, 0xef, 0xb6, 0x87, 0x2b, 0xd6, 0xa8, 0xa4, 0xc4, 0x8f,
	0xae, 0x3f, 0xf8, 0xd9, 0x0c, 0xf0, 0x42, 0x12, 0x07, 0xa6, 0x59, 0x9e, 0x9b, 0x66, 0xec, 0x2b,
	0xda, 0x19, 0xbc, 0x94, 0x59, 0x4e, 0x8a, 0xbe, 0xa7, 0x0c, 0x19, 0x35, 0x79, 0xb5, 0x8b, 0xe2,
	0xab, 0xfb, 0x6d, 0x22, 0xc7, 0x48, 0xc2, 0x31, 0x92, 0x32, 0x1c, 0xa3, 0x38, 0xfe, 0x29, 0x5b,
	0x7a, 0x70, 0x08, 0xd7, 0xbb, 0x28, 0x3e, 0x2b, 0x04, 0x34, 0xa8, 0x8b, 0x67, 0xfc, 0x65, 0x86,
	0x3c, 0x83, 0x0d, 0xb7, 0x09, 0xc8, 0x97, 0xaa, 0xba, 0x36, 0xcf, 0xe0, 0xf5, 0xf1, 0x52, 0x4c,
	0x94, 0xb9, 0x9c, 0xf2, 0x0c, 0xb4, 0x64, 0xa6, 0x9a, 0x5e, 0xf9, 0x64, 0x0f, 0x86, 0x96, 0xfd,
	0x66, 0x17, 0xc5, 0xcb, 0x22, 0x20, 0xa5, 0xfc, 0x62, 0x9c, 0x4f, 0x7d, 0x6b, 0xe0, 0x46, 0x52,
	0x06, 0x0e, 0x5e, 0x69, 0x26, 0x84, 0xb7, 0x27, 0x8f, 0x98, 0xb7, 0x38, 0xd6, 0x79, 0x06, 0xb7,
	0xfc, 0x9e, 0x40, 0xbd, 0xe2, 0x7c, 0x1f, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xf8, 0xf5, 0xd6,
	0xcc, 0x97, 0x02, 0x00, 0x00,
}