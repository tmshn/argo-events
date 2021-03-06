// Code generated by protoc-gen-go. DO NOT EDIT.
// source: eventing.proto

package gateways

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

//*
// Represents an event source
type EventSource struct {
	// ID of the event source. internally generated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The event source name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The event source configuration value.
	Value []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	// Type of the event source
	Type                 string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventSource) Reset()         { *m = EventSource{} }
func (m *EventSource) String() string { return proto.CompactTextString(m) }
func (*EventSource) ProtoMessage()    {}
func (*EventSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_2abcc01b0da84106, []int{0}
}

func (m *EventSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSource.Unmarshal(m, b)
}
func (m *EventSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSource.Marshal(b, m, deterministic)
}
func (m *EventSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSource.Merge(m, src)
}
func (m *EventSource) XXX_Size() int {
	return xxx_messageInfo_EventSource.Size(m)
}
func (m *EventSource) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSource.DiscardUnknown(m)
}

var xxx_messageInfo_EventSource proto.InternalMessageInfo

func (m *EventSource) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EventSource) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventSource) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *EventSource) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

//*
// Represents an event
type Event struct {
	// The event source name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The event payload.
	Payload              []byte   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_2abcc01b0da84106, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

//*
// Represents if an event source is valid or not
type ValidEventSource struct {
	// whether event source is valid
	IsValid bool `protobuf:"varint,1,opt,name=isValid,proto3" json:"isValid,omitempty"`
	// reason if an event source is invalid
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidEventSource) Reset()         { *m = ValidEventSource{} }
func (m *ValidEventSource) String() string { return proto.CompactTextString(m) }
func (*ValidEventSource) ProtoMessage()    {}
func (*ValidEventSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_2abcc01b0da84106, []int{2}
}

func (m *ValidEventSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidEventSource.Unmarshal(m, b)
}
func (m *ValidEventSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidEventSource.Marshal(b, m, deterministic)
}
func (m *ValidEventSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidEventSource.Merge(m, src)
}
func (m *ValidEventSource) XXX_Size() int {
	return xxx_messageInfo_ValidEventSource.Size(m)
}
func (m *ValidEventSource) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidEventSource.DiscardUnknown(m)
}

var xxx_messageInfo_ValidEventSource proto.InternalMessageInfo

func (m *ValidEventSource) GetIsValid() bool {
	if m != nil {
		return m.IsValid
	}
	return false
}

func (m *ValidEventSource) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func init() {
	proto.RegisterType((*EventSource)(nil), "gateways.EventSource")
	proto.RegisterType((*Event)(nil), "gateways.Event")
	proto.RegisterType((*ValidEventSource)(nil), "gateways.ValidEventSource")
}

func init() { proto.RegisterFile("eventing.proto", fileDescriptor_2abcc01b0da84106) }

var fileDescriptor_2abcc01b0da84106 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x86, 0xd9, 0xd8, 0x8f, 0x38, 0x96, 0x5a, 0xc6, 0x0f, 0x96, 0x9e, 0x4a, 0x4e, 0x3d, 0x05,
	0x51, 0xbc, 0x79, 0xb4, 0xe0, 0x39, 0x05, 0x2f, 0x9e, 0x46, 0x33, 0x94, 0x85, 0xb8, 0x1b, 0x36,
	0xdb, 0x4a, 0xfe, 0x86, 0xbf, 0x58, 0x32, 0x4d, 0x74, 0xe9, 0xc9, 0xdb, 0xbc, 0xcf, 0x0e, 0xcf,
	0xce, 0x0c, 0xcc, 0xf9, 0xc0, 0x36, 0x18, 0xbb, 0xcb, 0x6b, 0xef, 0x82, 0xc3, 0x74, 0x47, 0x81,
	0xbf, 0xa8, 0x6d, 0xb2, 0x37, 0xb8, 0xd8, 0x74, 0x6f, 0x5b, 0xb7, 0xf7, 0x1f, 0x8c, 0x73, 0x48,
	0x4c, 0xa9, 0xd5, 0x4a, 0xad, 0xcf, 0x8b, 0xc4, 0x94, 0x88, 0x30, 0xb2, 0xf4, 0xc9, 0x3a, 0x11,
	0x22, 0x35, 0x5e, 0xc3, 0xf8, 0x40, 0xd5, 0x9e, 0xf5, 0xd9, 0x4a, 0xad, 0x67, 0xc5, 0x31, 0x74,
	0x9d, 0xa1, 0xad, 0x59, 0x8f, 0x8e, 0x9d, 0x5d, 0x9d, 0x3d, 0xc2, 0x58, 0xe4, 0xbf, 0x1a, 0x15,
	0x69, 0x34, 0x4c, 0x6b, 0x6a, 0x2b, 0x47, 0xa5, 0xd8, 0x67, 0xc5, 0x10, 0xb3, 0x67, 0x58, 0xbc,
	0x52, 0x65, 0xca, 0x78, 0x30, 0x0d, 0x53, 0xd3, 0x08, 0x15, 0x49, 0x5a, 0x0c, 0x11, 0x6f, 0x61,
	0xe2, 0x99, 0x1a, 0x67, 0xfb, 0x21, 0xfb, 0x74, 0xff, 0xad, 0x20, 0xdd, 0xf4, 0x6b, 0xe3, 0x13,
	0x2c, 0xb6, 0x81, 0x7c, 0x88, 0x95, 0x37, 0xf9, 0x70, 0x85, 0x3c, 0xc2, 0xcb, 0xcb, 0x13, 0x7c,
	0xa7, 0xf0, 0x05, 0xae, 0xe4, 0x2f, 0x0a, 0xfc, 0x0f, 0xc1, 0xf2, 0x0f, 0x9f, 0xae, 0xf1, 0x3e,
	0x91, 0xfb, 0x3f, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xc0, 0xda, 0xf1, 0x91, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EventingClient is the client API for Eventing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventingClient interface {
	// StartEventSource starts an event source and returns stream of events.
	StartEventSource(ctx context.Context, in *EventSource, opts ...grpc.CallOption) (Eventing_StartEventSourceClient, error)
	// ValidateEventSource validates an event source.
	ValidateEventSource(ctx context.Context, in *EventSource, opts ...grpc.CallOption) (*ValidEventSource, error)
}

type eventingClient struct {
	cc *grpc.ClientConn
}

func NewEventingClient(cc *grpc.ClientConn) EventingClient {
	return &eventingClient{cc}
}

func (c *eventingClient) StartEventSource(ctx context.Context, in *EventSource, opts ...grpc.CallOption) (Eventing_StartEventSourceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Eventing_serviceDesc.Streams[0], "/gateways.Eventing/StartEventSource", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventingStartEventSourceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Eventing_StartEventSourceClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventingStartEventSourceClient struct {
	grpc.ClientStream
}

func (x *eventingStartEventSourceClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *eventingClient) ValidateEventSource(ctx context.Context, in *EventSource, opts ...grpc.CallOption) (*ValidEventSource, error) {
	out := new(ValidEventSource)
	err := c.cc.Invoke(ctx, "/gateways.Eventing/ValidateEventSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventingServer is the server API for Eventing service.
type EventingServer interface {
	// StartEventSource starts an event source and returns stream of events.
	StartEventSource(*EventSource, Eventing_StartEventSourceServer) error
	// ValidateEventSource validates an event source.
	ValidateEventSource(context.Context, *EventSource) (*ValidEventSource, error)
}

func RegisterEventingServer(s *grpc.Server, srv EventingServer) {
	s.RegisterService(&_Eventing_serviceDesc, srv)
}

func _Eventing_StartEventSource_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventSource)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventingServer).StartEventSource(m, &eventingStartEventSourceServer{stream})
}

type Eventing_StartEventSourceServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type eventingStartEventSourceServer struct {
	grpc.ServerStream
}

func (x *eventingStartEventSourceServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _Eventing_ValidateEventSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSource)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventingServer).ValidateEventSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateways.Eventing/ValidateEventSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventingServer).ValidateEventSource(ctx, req.(*EventSource))
	}
	return interceptor(ctx, in, info, handler)
}

var _Eventing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gateways.Eventing",
	HandlerType: (*EventingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateEventSource",
			Handler:    _Eventing_ValidateEventSource_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartEventSource",
			Handler:       _Eventing_StartEventSource_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "eventing.proto",
}
