// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wordcount.proto

package wordcount

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WordCount struct {
	Words                []string         `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	Count                map[string]int64 `protobuf:"bytes,2,rep,name=count,proto3" json:"count,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *WordCount) Reset()         { *m = WordCount{} }
func (m *WordCount) String() string { return proto.CompactTextString(m) }
func (*WordCount) ProtoMessage()    {}
func (*WordCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_wordcount_0b04b05e0de74145, []int{0}
}
func (m *WordCount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WordCount.Unmarshal(m, b)
}
func (m *WordCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WordCount.Marshal(b, m, deterministic)
}
func (dst *WordCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WordCount.Merge(dst, src)
}
func (m *WordCount) XXX_Size() int {
	return xxx_messageInfo_WordCount.Size(m)
}
func (m *WordCount) XXX_DiscardUnknown() {
	xxx_messageInfo_WordCount.DiscardUnknown(m)
}

var xxx_messageInfo_WordCount proto.InternalMessageInfo

func (m *WordCount) GetWords() []string {
	if m != nil {
		return m.Words
	}
	return nil
}

func (m *WordCount) GetCount() map[string]int64 {
	if m != nil {
		return m.Count
	}
	return nil
}

func init() {
	proto.RegisterType((*WordCount)(nil), "wordcount.WordCount")
	proto.RegisterMapType((map[string]int64)(nil), "wordcount.WordCount.CountEntry")
}

func init() { proto.RegisterFile("wordcount.proto", fileDescriptor_wordcount_0b04b05e0de74145) }

var fileDescriptor_wordcount_0b04b05e0de74145 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xcf, 0x2f, 0x4a,
	0x49, 0xce, 0x2f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x0b, 0x28,
	0x4d, 0x62, 0xe4, 0xe2, 0x0c, 0xcf, 0x2f, 0x4a, 0x71, 0x06, 0xf1, 0x84, 0x44, 0xb8, 0x58, 0x41,
	0x52, 0xc5, 0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90, 0x29, 0x17, 0x2b, 0x58,
	0xb1, 0x04, 0x93, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0xbc, 0x1e, 0xc2, 0x3c, 0xb8, 0x56, 0x3d, 0x30,
	0xe9, 0x9a, 0x57, 0x52, 0x54, 0x19, 0x04, 0x51, 0x2d, 0x65, 0xc1, 0xc5, 0x85, 0x10, 0x14, 0x12,
	0xe0, 0x62, 0xce, 0x4e, 0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x41, 0x96,
	0x95, 0x25, 0xe6, 0x94, 0xa6, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0x30, 0x07, 0x41, 0x38, 0x56, 0x4c,
	0x16, 0x8c, 0x49, 0x6c, 0x60, 0x67, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa6, 0x67, 0x25,
	0x0a, 0xb9, 0x00, 0x00, 0x00,
}
