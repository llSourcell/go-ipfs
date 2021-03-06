// Code generated by protoc-gen-go.
// source: encode.proto
// DO NOT EDIT!

/*
Package crypto is a generated protocol buffer package.

It is generated from these files:
	encode.proto

It has these top-level messages:
	PBPublicKey
	PBPrivateKey
*/
package crypto

import proto "github.com/jbenet/go-ipfs/Godeps/_workspace/src/code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type KeyType int32

const (
	KeyType_RSA KeyType = 0
)

var KeyType_name = map[int32]string{
	0: "RSA",
}
var KeyType_value = map[string]int32{
	"RSA": 0,
}

func (x KeyType) Enum() *KeyType {
	p := new(KeyType)
	*p = x
	return p
}
func (x KeyType) String() string {
	return proto.EnumName(KeyType_name, int32(x))
}
func (x *KeyType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(KeyType_value, data, "KeyType")
	if err != nil {
		return err
	}
	*x = KeyType(value)
	return nil
}

type PBPublicKey struct {
	Type             *KeyType `protobuf:"varint,1,req,enum=crypto.KeyType" json:"Type,omitempty"`
	Data             []byte   `protobuf:"bytes,2,req" json:"Data,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PBPublicKey) Reset()         { *m = PBPublicKey{} }
func (m *PBPublicKey) String() string { return proto.CompactTextString(m) }
func (*PBPublicKey) ProtoMessage()    {}

func (m *PBPublicKey) GetType() KeyType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return KeyType_RSA
}

func (m *PBPublicKey) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type PBPrivateKey struct {
	Type             *KeyType `protobuf:"varint,1,req,enum=crypto.KeyType" json:"Type,omitempty"`
	Data             []byte   `protobuf:"bytes,2,req" json:"Data,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PBPrivateKey) Reset()         { *m = PBPrivateKey{} }
func (m *PBPrivateKey) String() string { return proto.CompactTextString(m) }
func (*PBPrivateKey) ProtoMessage()    {}

func (m *PBPrivateKey) GetType() KeyType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return KeyType_RSA
}

func (m *PBPrivateKey) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("crypto.KeyType", KeyType_name, KeyType_value)
}
