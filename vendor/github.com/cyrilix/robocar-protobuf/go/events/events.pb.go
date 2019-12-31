// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events/events.proto

package events

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

type TypeObject int32

const (
	TypeObject_ANY  TypeObject = 0
	TypeObject_CAR  TypeObject = 1
	TypeObject_BUMP TypeObject = 2
	TypeObject_PLOT TypeObject = 3
)

var TypeObject_name = map[int32]string{
	0: "ANY",
	1: "CAR",
	2: "BUMP",
	3: "PLOT",
}

var TypeObject_value = map[string]int32{
	"ANY":  0,
	"CAR":  1,
	"BUMP": 2,
	"PLOT": 3,
}

func (x TypeObject) String() string {
	return proto.EnumName(TypeObject_name, int32(x))
}

func (TypeObject) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{0}
}

type FrameRef struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id                   string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FrameRef) Reset()         { *m = FrameRef{} }
func (m *FrameRef) String() string { return proto.CompactTextString(m) }
func (*FrameRef) ProtoMessage()    {}
func (*FrameRef) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{0}
}

func (m *FrameRef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FrameRef.Unmarshal(m, b)
}
func (m *FrameRef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FrameRef.Marshal(b, m, deterministic)
}
func (m *FrameRef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FrameRef.Merge(m, src)
}
func (m *FrameRef) XXX_Size() int {
	return xxx_messageInfo_FrameRef.Size(m)
}
func (m *FrameRef) XXX_DiscardUnknown() {
	xxx_messageInfo_FrameRef.DiscardUnknown(m)
}

var xxx_messageInfo_FrameRef proto.InternalMessageInfo

func (m *FrameRef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FrameRef) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type FrameMessage struct {
	Id                   *FrameRef `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Frame                []byte    `protobuf:"bytes,2,opt,name=frame,proto3" json:"frame,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *FrameMessage) Reset()         { *m = FrameMessage{} }
func (m *FrameMessage) String() string { return proto.CompactTextString(m) }
func (*FrameMessage) ProtoMessage()    {}
func (*FrameMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{1}
}

func (m *FrameMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FrameMessage.Unmarshal(m, b)
}
func (m *FrameMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FrameMessage.Marshal(b, m, deterministic)
}
func (m *FrameMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FrameMessage.Merge(m, src)
}
func (m *FrameMessage) XXX_Size() int {
	return xxx_messageInfo_FrameMessage.Size(m)
}
func (m *FrameMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_FrameMessage.DiscardUnknown(m)
}

var xxx_messageInfo_FrameMessage proto.InternalMessageInfo

func (m *FrameMessage) GetId() *FrameRef {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *FrameMessage) GetFrame() []byte {
	if m != nil {
		return m.Frame
	}
	return nil
}

type SteeringMessage struct {
	Steering             float32   `protobuf:"fixed32,1,opt,name=steering,proto3" json:"steering,omitempty"`
	Confidence           float32   `protobuf:"fixed32,2,opt,name=confidence,proto3" json:"confidence,omitempty"`
	FrameRef             *FrameRef `protobuf:"bytes,3,opt,name=frame_ref,json=frameRef,proto3" json:"frame_ref,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *SteeringMessage) Reset()         { *m = SteeringMessage{} }
func (m *SteeringMessage) String() string { return proto.CompactTextString(m) }
func (*SteeringMessage) ProtoMessage()    {}
func (*SteeringMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{2}
}

func (m *SteeringMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SteeringMessage.Unmarshal(m, b)
}
func (m *SteeringMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SteeringMessage.Marshal(b, m, deterministic)
}
func (m *SteeringMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SteeringMessage.Merge(m, src)
}
func (m *SteeringMessage) XXX_Size() int {
	return xxx_messageInfo_SteeringMessage.Size(m)
}
func (m *SteeringMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SteeringMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SteeringMessage proto.InternalMessageInfo

func (m *SteeringMessage) GetSteering() float32 {
	if m != nil {
		return m.Steering
	}
	return 0
}

func (m *SteeringMessage) GetConfidence() float32 {
	if m != nil {
		return m.Confidence
	}
	return 0
}

func (m *SteeringMessage) GetFrameRef() *FrameRef {
	if m != nil {
		return m.FrameRef
	}
	return nil
}

type ThrottleMessage struct {
	Throttle             float32   `protobuf:"fixed32,1,opt,name=throttle,proto3" json:"throttle,omitempty"`
	Confidence           float32   `protobuf:"fixed32,2,opt,name=confidence,proto3" json:"confidence,omitempty"`
	FrameRef             *FrameRef `protobuf:"bytes,3,opt,name=frame_ref,json=frameRef,proto3" json:"frame_ref,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ThrottleMessage) Reset()         { *m = ThrottleMessage{} }
func (m *ThrottleMessage) String() string { return proto.CompactTextString(m) }
func (*ThrottleMessage) ProtoMessage()    {}
func (*ThrottleMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{3}
}

func (m *ThrottleMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ThrottleMessage.Unmarshal(m, b)
}
func (m *ThrottleMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ThrottleMessage.Marshal(b, m, deterministic)
}
func (m *ThrottleMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThrottleMessage.Merge(m, src)
}
func (m *ThrottleMessage) XXX_Size() int {
	return xxx_messageInfo_ThrottleMessage.Size(m)
}
func (m *ThrottleMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ThrottleMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ThrottleMessage proto.InternalMessageInfo

func (m *ThrottleMessage) GetThrottle() float32 {
	if m != nil {
		return m.Throttle
	}
	return 0
}

func (m *ThrottleMessage) GetConfidence() float32 {
	if m != nil {
		return m.Confidence
	}
	return 0
}

func (m *ThrottleMessage) GetFrameRef() *FrameRef {
	if m != nil {
		return m.FrameRef
	}
	return nil
}

type ObjectsMessage struct {
	Objects              []*Object `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
	FrameRef             *FrameRef `protobuf:"bytes,2,opt,name=frame_ref,json=frameRef,proto3" json:"frame_ref,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ObjectsMessage) Reset()         { *m = ObjectsMessage{} }
func (m *ObjectsMessage) String() string { return proto.CompactTextString(m) }
func (*ObjectsMessage) ProtoMessage()    {}
func (*ObjectsMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{4}
}

func (m *ObjectsMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObjectsMessage.Unmarshal(m, b)
}
func (m *ObjectsMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObjectsMessage.Marshal(b, m, deterministic)
}
func (m *ObjectsMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObjectsMessage.Merge(m, src)
}
func (m *ObjectsMessage) XXX_Size() int {
	return xxx_messageInfo_ObjectsMessage.Size(m)
}
func (m *ObjectsMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ObjectsMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ObjectsMessage proto.InternalMessageInfo

func (m *ObjectsMessage) GetObjects() []*Object {
	if m != nil {
		return m.Objects
	}
	return nil
}

func (m *ObjectsMessage) GetFrameRef() *FrameRef {
	if m != nil {
		return m.FrameRef
	}
	return nil
}

type Object struct {
	Type                 TypeObject `protobuf:"varint,1,opt,name=type,proto3,enum=robocar.events.TypeObject" json:"type,omitempty"`
	LLeft                int32      `protobuf:"varint,2,opt,name=lLeft,proto3" json:"lLeft,omitempty"`
	Up                   int32      `protobuf:"varint,3,opt,name=up,proto3" json:"up,omitempty"`
	Right                int32      `protobuf:"varint,4,opt,name=right,proto3" json:"right,omitempty"`
	Bottom               int32      `protobuf:"varint,5,opt,name=bottom,proto3" json:"bottom,omitempty"`
	Confidence           float32    `protobuf:"fixed32,6,opt,name=confidence,proto3" json:"confidence,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Object) Reset()         { *m = Object{} }
func (m *Object) String() string { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()    {}
func (*Object) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ec31f2d2a3db598, []int{5}
}

func (m *Object) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Object.Unmarshal(m, b)
}
func (m *Object) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Object.Marshal(b, m, deterministic)
}
func (m *Object) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object.Merge(m, src)
}
func (m *Object) XXX_Size() int {
	return xxx_messageInfo_Object.Size(m)
}
func (m *Object) XXX_DiscardUnknown() {
	xxx_messageInfo_Object.DiscardUnknown(m)
}

var xxx_messageInfo_Object proto.InternalMessageInfo

func (m *Object) GetType() TypeObject {
	if m != nil {
		return m.Type
	}
	return TypeObject_ANY
}

func (m *Object) GetLLeft() int32 {
	if m != nil {
		return m.LLeft
	}
	return 0
}

func (m *Object) GetUp() int32 {
	if m != nil {
		return m.Up
	}
	return 0
}

func (m *Object) GetRight() int32 {
	if m != nil {
		return m.Right
	}
	return 0
}

func (m *Object) GetBottom() int32 {
	if m != nil {
		return m.Bottom
	}
	return 0
}

func (m *Object) GetConfidence() float32 {
	if m != nil {
		return m.Confidence
	}
	return 0
}

func init() {
	proto.RegisterEnum("robocar.events.TypeObject", TypeObject_name, TypeObject_value)
	proto.RegisterType((*FrameRef)(nil), "robocar.events.FrameRef")
	proto.RegisterType((*FrameMessage)(nil), "robocar.events.FrameMessage")
	proto.RegisterType((*SteeringMessage)(nil), "robocar.events.SteeringMessage")
	proto.RegisterType((*ThrottleMessage)(nil), "robocar.events.ThrottleMessage")
	proto.RegisterType((*ObjectsMessage)(nil), "robocar.events.ObjectsMessage")
	proto.RegisterType((*Object)(nil), "robocar.events.Object")
}

func init() { proto.RegisterFile("events/events.proto", fileDescriptor_8ec31f2d2a3db598) }

var fileDescriptor_8ec31f2d2a3db598 = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x4d, 0x6f, 0x9b, 0x40,
	0x10, 0x2d, 0xcb, 0x87, 0xf1, 0xd8, 0xa2, 0x68, 0x5b, 0x59, 0xc8, 0x87, 0xca, 0xf2, 0xc9, 0xea,
	0x81, 0x56, 0x54, 0xfd, 0x01, 0x76, 0xa5, 0x9e, 0xfc, 0xa5, 0xad, 0x7b, 0x48, 0x2e, 0x11, 0xe0,
	0xc5, 0x26, 0xb2, 0x59, 0xb4, 0xac, 0x23, 0xf9, 0x9e, 0x9f, 0x93, 0x1f, 0x19, 0xed, 0x2e, 0x24,
	0x36, 0x97, 0xe4, 0x92, 0x13, 0xf3, 0x1e, 0x6f, 0xde, 0xbc, 0x59, 0x58, 0xf8, 0x42, 0x1f, 0x68,
	0x21, 0xaa, 0x1f, 0xfa, 0x11, 0x96, 0x9c, 0x09, 0x86, 0x3d, 0xce, 0x12, 0x96, 0xc6, 0x3c, 0xd4,
	0xec, 0x38, 0x04, 0xf7, 0x2f, 0x8f, 0x8f, 0x94, 0xd0, 0x0c, 0x63, 0xb0, 0x8a, 0xf8, 0x48, 0x03,
	0x63, 0x64, 0x4c, 0xba, 0x44, 0xd5, 0xd8, 0x03, 0x94, 0x6f, 0x03, 0xa4, 0x18, 0x94, 0x6f, 0xc7,
	0x4b, 0xe8, 0x2b, 0xfd, 0x82, 0x56, 0x55, 0xbc, 0xa3, 0x78, 0xa2, 0xde, 0xcb, 0x8e, 0x5e, 0x14,
	0x84, 0xd7, 0xe6, 0x61, 0xe3, 0x2c, 0x3b, 0xf1, 0x57, 0xb0, 0x33, 0x89, 0x95, 0x59, 0x9f, 0x68,
	0x30, 0x7e, 0x34, 0xe0, 0xf3, 0x3f, 0x41, 0x29, 0xcf, 0x8b, 0x5d, 0xe3, 0x39, 0x04, 0xb7, 0xaa,
	0x29, 0xe5, 0x8c, 0xc8, 0x0b, 0xc6, 0xdf, 0x00, 0x52, 0x56, 0x64, 0xf9, 0x96, 0x16, 0xa9, 0xb6,
	0x42, 0xe4, 0x82, 0xc1, 0xbf, 0xa1, 0xab, 0x8c, 0xef, 0x38, 0xcd, 0x02, 0xf3, 0x8d, 0x58, 0x6e,
	0x56, 0x57, 0x2a, 0xc6, 0x66, 0xcf, 0x99, 0x10, 0x07, 0x7a, 0x11, 0x43, 0xd4, 0x54, 0x13, 0xa3,
	0xc1, 0x1f, 0x15, 0xe3, 0x0c, 0xde, 0x2a, 0xb9, 0xa7, 0xa9, 0xa8, 0x9a, 0x10, 0x3f, 0xa1, 0xc3,
	0x34, 0x13, 0x18, 0x23, 0x73, 0xd2, 0x8b, 0x06, 0x6d, 0x1b, 0xdd, 0x40, 0x1a, 0xd9, 0xf5, 0x68,
	0xf4, 0xee, 0xd1, 0x4f, 0x06, 0x38, 0xda, 0x0a, 0x87, 0x60, 0x89, 0x73, 0xa9, 0x97, 0xf6, 0xa2,
	0x61, 0xbb, 0x79, 0x73, 0x2e, 0x69, 0x3d, 0x54, 0xe9, 0xe4, 0x97, 0x3d, 0xcc, 0x69, 0x26, 0xd4,
	0x34, 0x9b, 0x68, 0x20, 0xff, 0x9c, 0x53, 0xa9, 0x76, 0xb7, 0x09, 0x3a, 0x95, 0x52, 0xc5, 0xf3,
	0xdd, 0x5e, 0x04, 0x96, 0x56, 0x29, 0x80, 0x07, 0xe0, 0x24, 0x4c, 0x08, 0x76, 0x0c, 0x6c, 0x45,
	0xd7, 0xa8, 0x75, 0xc0, 0x4e, 0xfb, 0x80, 0xbf, 0x47, 0x00, 0xaf, 0x39, 0x70, 0x07, 0xcc, 0xe9,
	0xf2, 0xc6, 0xff, 0x24, 0x8b, 0x3f, 0x53, 0xe2, 0x1b, 0xd8, 0x05, 0x6b, 0xf6, 0x7f, 0xb1, 0xf6,
	0x91, 0xac, 0xd6, 0xf3, 0xd5, 0xc6, 0x37, 0x67, 0xee, 0xad, 0xa3, 0x57, 0x48, 0x1c, 0x75, 0x19,
	0x7e, 0x3d, 0x07, 0x00, 0x00, 0xff, 0xff, 0xb6, 0x83, 0x51, 0xbe, 0x23, 0x03, 0x00, 0x00,
}
