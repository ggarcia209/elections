// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmte_tx_data.proto

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

type CmteEntry struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Total                float32  `protobuf:"fixed32,2,opt,name=Total,proto3" json:"Total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CmteEntry) Reset()         { *m = CmteEntry{} }
func (m *CmteEntry) String() string { return proto.CompactTextString(m) }
func (*CmteEntry) ProtoMessage()    {}
func (*CmteEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_e66b7cd10fa5e378, []int{0}
}

func (m *CmteEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmteEntry.Unmarshal(m, b)
}
func (m *CmteEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmteEntry.Marshal(b, m, deterministic)
}
func (m *CmteEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmteEntry.Merge(m, src)
}
func (m *CmteEntry) XXX_Size() int {
	return xxx_messageInfo_CmteEntry.Size(m)
}
func (m *CmteEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_CmteEntry.DiscardUnknown(m)
}

var xxx_messageInfo_CmteEntry proto.InternalMessageInfo

func (m *CmteEntry) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *CmteEntry) GetTotal() float32 {
	if m != nil {
		return m.Total
	}
	return 0
}

type CmteTxData struct {
	CmteID                         string             `protobuf:"bytes,1,opt,name=CmteID,proto3" json:"CmteID,omitempty"`
	CandID                         string             `protobuf:"bytes,2,opt,name=CandID,proto3" json:"CandID,omitempty"`
	Party                          string             `protobuf:"bytes,3,opt,name=Party,proto3" json:"Party,omitempty"`
	ContributionsInAmt             float32            `protobuf:"fixed32,4,opt,name=ContributionsInAmt,proto3" json:"ContributionsInAmt,omitempty"`
	ContributionsInTxs             float32            `protobuf:"fixed32,5,opt,name=ContributionsInTxs,proto3" json:"ContributionsInTxs,omitempty"`
	AvgContributionIn              float32            `protobuf:"fixed32,6,opt,name=AvgContributionIn,proto3" json:"AvgContributionIn,omitempty"`
	OtherReceiptsInAmt             float32            `protobuf:"fixed32,7,opt,name=OtherReceiptsInAmt,proto3" json:"OtherReceiptsInAmt,omitempty"`
	OtherReceiptsInTxs             float32            `protobuf:"fixed32,8,opt,name=OtherReceiptsInTxs,proto3" json:"OtherReceiptsInTxs,omitempty"`
	AvgOtherIn                     float32            `protobuf:"fixed32,9,opt,name=AvgOtherIn,proto3" json:"AvgOtherIn,omitempty"`
	TotalIncomingAmt               float32            `protobuf:"fixed32,10,opt,name=TotalIncomingAmt,proto3" json:"TotalIncomingAmt,omitempty"`
	TotalIncomingTxs               float32            `protobuf:"fixed32,11,opt,name=TotalIncomingTxs,proto3" json:"TotalIncomingTxs,omitempty"`
	AvgIncoming                    float32            `protobuf:"fixed32,12,opt,name=AvgIncoming,proto3" json:"AvgIncoming,omitempty"`
	TransfersAmt                   float32            `protobuf:"fixed32,13,opt,name=TransfersAmt,proto3" json:"TransfersAmt,omitempty"`
	TransfersTxs                   float32            `protobuf:"fixed32,14,opt,name=TransfersTxs,proto3" json:"TransfersTxs,omitempty"`
	AvgTransfer                    float32            `protobuf:"fixed32,15,opt,name=AvgTransfer,proto3" json:"AvgTransfer,omitempty"`
	ExpendituresAmt                float32            `protobuf:"fixed32,16,opt,name=ExpendituresAmt,proto3" json:"ExpendituresAmt,omitempty"`
	ExpendituresTxs                float32            `protobuf:"fixed32,17,opt,name=ExpendituresTxs,proto3" json:"ExpendituresTxs,omitempty"`
	AvgExpenditure                 float32            `protobuf:"fixed32,18,opt,name=AvgExpenditure,proto3" json:"AvgExpenditure,omitempty"`
	TotalOutgoingAmt               float32            `protobuf:"fixed32,19,opt,name=TotalOutgoingAmt,proto3" json:"TotalOutgoingAmt,omitempty"`
	TotalOutgoingTxs               float32            `protobuf:"fixed32,20,opt,name=TotalOutgoingTxs,proto3" json:"TotalOutgoingTxs,omitempty"`
	AvgOutgoing                    float32            `protobuf:"fixed32,21,opt,name=AvgOutgoing,proto3" json:"AvgOutgoing,omitempty"`
	NetBalance                     float32            `protobuf:"fixed32,22,opt,name=NetBalance,proto3" json:"NetBalance,omitempty"`
	TopIndvContributorsAmt         map[string]float32 `protobuf:"bytes,23,rep,name=TopIndvContributorsAmt,proto3" json:"TopIndvContributorsAmt,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopIndvContributorsTxs         map[string]float32 `protobuf:"bytes,24,rep,name=TopIndvContributorsTxs,proto3" json:"TopIndvContributorsTxs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopIndvContributorThreshold    []*CmteEntry       `protobuf:"bytes,25,rep,name=TopIndvContributorThreshold,proto3" json:"TopIndvContributorThreshold,omitempty"`
	TopCmteOrgContributorsAmt      map[string]float32 `protobuf:"bytes,26,rep,name=TopCmteOrgContributorsAmt,proto3" json:"TopCmteOrgContributorsAmt,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopCmteOrgContributorsTxs      map[string]float32 `protobuf:"bytes,27,rep,name=TopCmteOrgContributorsTxs,proto3" json:"TopCmteOrgContributorsTxs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopCmteOrgContributorThreshold []*CmteEntry       `protobuf:"bytes,28,rep,name=TopCmteOrgContributorThreshold,proto3" json:"TopCmteOrgContributorThreshold,omitempty"`
	TransferRecsAmt                map[string]float32 `protobuf:"bytes,29,rep,name=TransferRecsAmt,proto3" json:"TransferRecsAmt,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TransferRecsTxs                map[string]float32 `protobuf:"bytes,30,rep,name=TransferRecsTxs,proto3" json:"TransferRecsTxs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopExpRecipientsAmt            map[string]float32 `protobuf:"bytes,31,rep,name=TopExpRecipientsAmt,proto3" json:"TopExpRecipientsAmt,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopExpRecipientsTxs            map[string]float32 `protobuf:"bytes,32,rep,name=TopExpRecipientsTxs,proto3" json:"TopExpRecipientsTxs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	TopExpThreshold                []*CmteEntry       `protobuf:"bytes,33,rep,name=TopExpThreshold,proto3" json:"TopExpThreshold,omitempty"`
	XXX_NoUnkeyedLiteral           struct{}           `json:"-"`
	XXX_unrecognized               []byte             `json:"-"`
	XXX_sizecache                  int32              `json:"-"`
}

func (m *CmteTxData) Reset()         { *m = CmteTxData{} }
func (m *CmteTxData) String() string { return proto.CompactTextString(m) }
func (*CmteTxData) ProtoMessage()    {}
func (*CmteTxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_e66b7cd10fa5e378, []int{1}
}

func (m *CmteTxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmteTxData.Unmarshal(m, b)
}
func (m *CmteTxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmteTxData.Marshal(b, m, deterministic)
}
func (m *CmteTxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmteTxData.Merge(m, src)
}
func (m *CmteTxData) XXX_Size() int {
	return xxx_messageInfo_CmteTxData.Size(m)
}
func (m *CmteTxData) XXX_DiscardUnknown() {
	xxx_messageInfo_CmteTxData.DiscardUnknown(m)
}

var xxx_messageInfo_CmteTxData proto.InternalMessageInfo

func (m *CmteTxData) GetCmteID() string {
	if m != nil {
		return m.CmteID
	}
	return ""
}

func (m *CmteTxData) GetCandID() string {
	if m != nil {
		return m.CandID
	}
	return ""
}

func (m *CmteTxData) GetParty() string {
	if m != nil {
		return m.Party
	}
	return ""
}

func (m *CmteTxData) GetContributionsInAmt() float32 {
	if m != nil {
		return m.ContributionsInAmt
	}
	return 0
}

func (m *CmteTxData) GetContributionsInTxs() float32 {
	if m != nil {
		return m.ContributionsInTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgContributionIn() float32 {
	if m != nil {
		return m.AvgContributionIn
	}
	return 0
}

func (m *CmteTxData) GetOtherReceiptsInAmt() float32 {
	if m != nil {
		return m.OtherReceiptsInAmt
	}
	return 0
}

func (m *CmteTxData) GetOtherReceiptsInTxs() float32 {
	if m != nil {
		return m.OtherReceiptsInTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgOtherIn() float32 {
	if m != nil {
		return m.AvgOtherIn
	}
	return 0
}

func (m *CmteTxData) GetTotalIncomingAmt() float32 {
	if m != nil {
		return m.TotalIncomingAmt
	}
	return 0
}

func (m *CmteTxData) GetTotalIncomingTxs() float32 {
	if m != nil {
		return m.TotalIncomingTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgIncoming() float32 {
	if m != nil {
		return m.AvgIncoming
	}
	return 0
}

func (m *CmteTxData) GetTransfersAmt() float32 {
	if m != nil {
		return m.TransfersAmt
	}
	return 0
}

func (m *CmteTxData) GetTransfersTxs() float32 {
	if m != nil {
		return m.TransfersTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgTransfer() float32 {
	if m != nil {
		return m.AvgTransfer
	}
	return 0
}

func (m *CmteTxData) GetExpendituresAmt() float32 {
	if m != nil {
		return m.ExpendituresAmt
	}
	return 0
}

func (m *CmteTxData) GetExpendituresTxs() float32 {
	if m != nil {
		return m.ExpendituresTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgExpenditure() float32 {
	if m != nil {
		return m.AvgExpenditure
	}
	return 0
}

func (m *CmteTxData) GetTotalOutgoingAmt() float32 {
	if m != nil {
		return m.TotalOutgoingAmt
	}
	return 0
}

func (m *CmteTxData) GetTotalOutgoingTxs() float32 {
	if m != nil {
		return m.TotalOutgoingTxs
	}
	return 0
}

func (m *CmteTxData) GetAvgOutgoing() float32 {
	if m != nil {
		return m.AvgOutgoing
	}
	return 0
}

func (m *CmteTxData) GetNetBalance() float32 {
	if m != nil {
		return m.NetBalance
	}
	return 0
}

func (m *CmteTxData) GetTopIndvContributorsAmt() map[string]float32 {
	if m != nil {
		return m.TopIndvContributorsAmt
	}
	return nil
}

func (m *CmteTxData) GetTopIndvContributorsTxs() map[string]float32 {
	if m != nil {
		return m.TopIndvContributorsTxs
	}
	return nil
}

func (m *CmteTxData) GetTopIndvContributorThreshold() []*CmteEntry {
	if m != nil {
		return m.TopIndvContributorThreshold
	}
	return nil
}

func (m *CmteTxData) GetTopCmteOrgContributorsAmt() map[string]float32 {
	if m != nil {
		return m.TopCmteOrgContributorsAmt
	}
	return nil
}

func (m *CmteTxData) GetTopCmteOrgContributorsTxs() map[string]float32 {
	if m != nil {
		return m.TopCmteOrgContributorsTxs
	}
	return nil
}

func (m *CmteTxData) GetTopCmteOrgContributorThreshold() []*CmteEntry {
	if m != nil {
		return m.TopCmteOrgContributorThreshold
	}
	return nil
}

func (m *CmteTxData) GetTransferRecsAmt() map[string]float32 {
	if m != nil {
		return m.TransferRecsAmt
	}
	return nil
}

func (m *CmteTxData) GetTransferRecsTxs() map[string]float32 {
	if m != nil {
		return m.TransferRecsTxs
	}
	return nil
}

func (m *CmteTxData) GetTopExpRecipientsAmt() map[string]float32 {
	if m != nil {
		return m.TopExpRecipientsAmt
	}
	return nil
}

func (m *CmteTxData) GetTopExpRecipientsTxs() map[string]float32 {
	if m != nil {
		return m.TopExpRecipientsTxs
	}
	return nil
}

func (m *CmteTxData) GetTopExpThreshold() []*CmteEntry {
	if m != nil {
		return m.TopExpThreshold
	}
	return nil
}

func init() {
	proto.RegisterType((*CmteEntry)(nil), "protobuf.CmteEntry")
	proto.RegisterType((*CmteTxData)(nil), "protobuf.CmteTxData")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopCmteOrgContributorsAmtEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopCmteOrgContributorsTxsEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopExpRecipientsAmtEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopExpRecipientsTxsEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopIndvContributorsAmtEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TopIndvContributorsTxsEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TransferRecsAmtEntry")
	proto.RegisterMapType((map[string]float32)(nil), "protobuf.CmteTxData.TransferRecsTxsEntry")
}

func init() { proto.RegisterFile("cmte_tx_data.proto", fileDescriptor_e66b7cd10fa5e378) }

var fileDescriptor_e66b7cd10fa5e378 = []byte{
	// 684 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x55, 0x3b, 0x36, 0xb6, 0xbb, 0xb1, 0x6e, 0xee, 0x18, 0xde, 0x06, 0xa3, 0xec, 0x01, 0x15,
	0x04, 0x15, 0xb0, 0x17, 0x84, 0xc4, 0x43, 0xb7, 0x0e, 0x29, 0x12, 0x62, 0xa8, 0x84, 0x27, 0x1e,
	0xa6, 0xac, 0xf5, 0xd2, 0x88, 0xd6, 0x09, 0x89, 0x5b, 0x65, 0x1f, 0xc7, 0xbf, 0xa1, 0x7b, 0xd3,
	0xa4, 0x69, 0x62, 0x03, 0xed, 0x9e, 0x5a, 0x9f, 0x7b, 0xee, 0x39, 0xd7, 0xc7, 0xb6, 0x02, 0xac,
	0x37, 0x52, 0xe2, 0x4a, 0xc5, 0x57, 0x7d, 0x47, 0x39, 0xad, 0x20, 0xf4, 0x95, 0xcf, 0xd6, 0xe9,
	0xe7, 0x7a, 0x7c, 0x73, 0xf2, 0x16, 0x36, 0xce, 0x47, 0x4a, 0x5c, 0x48, 0x15, 0xde, 0xb2, 0x6d,
	0xa8, 0x5a, 0x1d, 0x5e, 0x69, 0x54, 0x9a, 0x1b, 0xdd, 0xaa, 0xd5, 0x61, 0x7b, 0xb0, 0x6a, 0xfb,
	0xca, 0x19, 0xf2, 0x6a, 0xa3, 0xd2, 0xac, 0x76, 0x93, 0xc5, 0xc9, 0xef, 0x3a, 0x00, 0xf6, 0xd8,
	0x71, 0xc7, 0x51, 0x0e, 0xdb, 0x87, 0x35, 0x5c, 0x65, 0x8d, 0xd3, 0x15, 0xe1, 0x8e, 0xec, 0x5b,
	0x1d, 0xea, 0x46, 0x9c, 0x56, 0x28, 0xfa, 0xd5, 0x09, 0xd5, 0x2d, 0x5f, 0x21, 0x38, 0x59, 0xb0,
	0x16, 0xb0, 0x73, 0x5f, 0xaa, 0xd0, 0xbb, 0x1e, 0x2b, 0xcf, 0x97, 0x91, 0x25, 0xdb, 0x23, 0xc5,
	0xef, 0x91, 0xaf, 0xa6, 0xa2, 0xe1, 0xdb, 0x71, 0xc4, 0x57, 0xb5, 0x7c, 0x3b, 0x8e, 0xd8, 0x2b,
	0xd8, 0x6d, 0x4f, 0xdc, 0x7c, 0xc1, 0x92, 0x7c, 0x8d, 0xe8, 0xe5, 0x02, 0xaa, 0x5f, 0xaa, 0x81,
	0x08, 0xbb, 0xa2, 0x27, 0xbc, 0x40, 0x4d, 0xa7, 0xb9, 0x9f, 0xa8, 0x97, 0x2b, 0x1a, 0x3e, 0x4e,
	0xb3, 0xae, 0xe5, 0xe3, 0x34, 0xc7, 0x00, 0xed, 0x89, 0x4b, 0x05, 0x4b, 0xf2, 0x0d, 0xe2, 0xe5,
	0x10, 0xf6, 0x12, 0x76, 0x28, 0x6b, 0x4b, 0xf6, 0xfc, 0x91, 0x27, 0x5d, 0x74, 0x07, 0x62, 0x95,
	0xf0, 0x12, 0x17, 0x9d, 0x37, 0x35, 0x5c, 0xf4, 0x6d, 0xc0, 0x66, 0x7b, 0xe2, 0xa6, 0x08, 0xdf,
	0x22, 0x5a, 0x1e, 0x62, 0x27, 0xb0, 0x65, 0x87, 0x8e, 0x8c, 0x6e, 0x44, 0x18, 0xa1, 0xeb, 0x03,
	0xa2, 0xcc, 0x61, 0x73, 0x1c, 0x74, 0xdb, 0x2e, 0x70, 0x66, 0x4e, 0x29, 0xc4, 0x6b, 0x99, 0x53,
	0x0a, 0xb1, 0x26, 0xd4, 0x2e, 0xe2, 0x40, 0xc8, 0xbe, 0xa7, 0xc6, 0xa1, 0x20, 0xb3, 0x1d, 0x62,
	0x15, 0xe1, 0x22, 0x13, 0x2d, 0x77, 0xcb, 0x4c, 0x74, 0x7d, 0x0e, 0xdb, 0xed, 0x89, 0x9b, 0x43,
	0x39, 0x23, 0x62, 0x01, 0xcd, 0x32, 0xbb, 0x1c, 0x2b, 0xd7, 0x9f, 0xe6, 0x5b, 0xcf, 0x65, 0x96,
	0xc3, 0x4b, 0x5c, 0xb4, 0xdf, 0xd3, 0x70, 0x67, 0xbb, 0x4e, 0x11, 0xfe, 0x30, 0xdb, 0x75, 0x0a,
	0xe1, 0xc9, 0x7f, 0x11, 0xea, 0xcc, 0x19, 0x3a, 0xb2, 0x27, 0xf8, 0x7e, 0x72, 0xf2, 0x33, 0x84,
	0x0d, 0x60, 0xdf, 0xf6, 0x03, 0x4b, 0xf6, 0x27, 0xd9, 0x95, 0xf4, 0x93, 0x93, 0x78, 0xd4, 0x58,
	0x69, 0x6e, 0xbe, 0x7b, 0xd3, 0x4a, 0x9f, 0x6e, 0x6b, 0xf6, 0x06, 0x5b, 0xfa, 0x16, 0x7a, 0xd4,
	0x5d, 0x83, 0x9e, 0xc1, 0x09, 0x77, 0xc7, 0x17, 0x73, 0xb2, 0xe3, 0xc8, 0xec, 0x84, 0xa9, 0x7c,
	0x87, 0xa3, 0x72, 0xc5, 0x1e, 0x84, 0x22, 0x1a, 0xf8, 0xc3, 0x3e, 0x3f, 0x20, 0xbb, 0xfa, 0xbc,
	0x5d, 0xa2, 0xf8, 0xb7, 0x3e, 0xf6, 0x0b, 0x0e, 0x6c, 0x3f, 0x40, 0xf2, 0x65, 0xe8, 0x16, 0xd3,
	0x3a, 0x24, 0xd1, 0x53, 0xd3, 0x1e, 0xf4, 0x5d, 0x89, 0xa9, 0x59, 0xd5, 0x6c, 0x89, 0xb1, 0x1d,
	0x2d, 0x6c, 0x99, 0x25, 0x67, 0x56, 0x65, 0x3f, 0xe0, 0x58, 0x5b, 0x9c, 0xe5, 0xf7, 0xd8, 0x9c,
	0xdf, 0x3f, 0x5a, 0xd9, 0x37, 0xa8, 0xa5, 0xef, 0xb1, 0x2b, 0x7a, 0x14, 0xdc, 0x13, 0x52, 0x7b,
	0xa1, 0xdf, 0xc5, 0x3c, 0x37, 0xf1, 0x28, 0x2a, 0x14, 0x45, 0x31, 0x9a, 0xe3, 0xff, 0x14, 0xcd,
	0x02, 0x29, 0x2a, 0xb0, 0x2b, 0xa8, 0xdb, 0x7e, 0x70, 0x11, 0x07, 0x5d, 0xd1, 0xf3, 0x02, 0x4f,
	0x48, 0x45, 0xd3, 0x3e, 0x25, 0xe1, 0xd7, 0xa6, 0xcc, 0x8b, 0xfc, 0x44, 0x5c, 0xa7, 0xa4, 0x33,
	0xc0, 0xc9, 0x1b, 0x0b, 0x18, 0x64, 0xd3, 0xeb, 0x94, 0xd8, 0x47, 0xa8, 0x25, 0xf0, 0xec, 0xe4,
	0x9e, 0x99, 0x4f, 0xae, 0xc8, 0x3d, 0xb4, 0x74, 0x8f, 0x28, 0xdb, 0x13, 0xdb, 0x81, 0x95, 0x9f,
	0xe2, 0x76, 0xfa, 0x09, 0xc6, 0xbf, 0xf8, 0x9d, 0x9d, 0x38, 0xc3, 0xb1, 0x48, 0x3f, 0xde, 0xb4,
	0xf8, 0x50, 0x7d, 0x5f, 0x31, 0x48, 0xa5, 0xd3, 0x2f, 0x24, 0xf5, 0xd9, 0x70, 0x3b, 0x97, 0x1b,
	0xcc, 0xa8, 0xb6, 0xd4, 0x6c, 0x67, 0xb0, 0xa7, 0xbb, 0xb0, 0x77, 0xd1, 0x58, 0x6a, 0x8e, 0x4f,
	0xc0, 0x4d, 0x57, 0xf1, 0xae, 0x3a, 0xcb, 0xcc, 0x73, 0xbd, 0x46, 0xd7, 0xed, 0xf4, 0x4f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xb3, 0x75, 0x50, 0x7a, 0x19, 0x0a, 0x00, 0x00,
}