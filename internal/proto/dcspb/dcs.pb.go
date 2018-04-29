// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dcs.proto

/*
Package dcspb is a generated protocol buffer package.

It is generated from these files:
	dcs.proto

It has these top-level messages:
	SearchRequest
	Error
	Progress
	Pagination
	Event
*/
package dcspb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import sourcebackendpb "github.com/Debian/dcs/internal/proto/sourcebackendpb"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Error_ErrorType int32

const (
	Error_CANCELLED           Error_ErrorType = 0
	Error_BACKEND_UNAVAILABLE Error_ErrorType = 1
	Error_FAILED              Error_ErrorType = 2
	Error_INVALID_QUERY       Error_ErrorType = 3
)

var Error_ErrorType_name = map[int32]string{
	0: "CANCELLED",
	1: "BACKEND_UNAVAILABLE",
	2: "FAILED",
	3: "INVALID_QUERY",
}
var Error_ErrorType_value = map[string]int32{
	"CANCELLED":           0,
	"BACKEND_UNAVAILABLE": 1,
	"FAILED":              2,
	"INVALID_QUERY":       3,
}

func (x Error_ErrorType) String() string {
	return proto.EnumName(Error_ErrorType_name, int32(x))
}
func (Error_ErrorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type Event_Type int32

const (
	Event_ERROR      Event_Type = 0
	Event_PROGRESS   Event_Type = 1
	Event_MATCH      Event_Type = 2
	Event_PAGINATION Event_Type = 3
	Event_DONE       Event_Type = 4
)

var Event_Type_name = map[int32]string{
	0: "ERROR",
	1: "PROGRESS",
	2: "MATCH",
	3: "PAGINATION",
	4: "DONE",
}
var Event_Type_value = map[string]int32{
	"ERROR":      0,
	"PROGRESS":   1,
	"MATCH":      2,
	"PAGINATION": 3,
	"DONE":       4,
}

func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}
func (Event_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

type SearchRequest struct {
	Query string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SearchRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type Error struct {
	Type    Error_ErrorType `protobuf:"varint,1,opt,name=type,enum=dcspb.Error_ErrorType" json:"type,omitempty"`
	Message string          `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Error) GetType() Error_ErrorType {
	if m != nil {
		return m.Type
	}
	return Error_CANCELLED
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Progress struct {
	QueryId        string `protobuf:"bytes,1,opt,name=query_id,json=queryId" json:"query_id,omitempty"`
	FilesProcessed int64  `protobuf:"varint,2,opt,name=files_processed,json=filesProcessed" json:"files_processed,omitempty"`
	FilesTotal     int64  `protobuf:"varint,3,opt,name=files_total,json=filesTotal" json:"files_total,omitempty"`
	Results        int64  `protobuf:"varint,4,opt,name=results" json:"results,omitempty"`
}

func (m *Progress) Reset()                    { *m = Progress{} }
func (m *Progress) String() string            { return proto.CompactTextString(m) }
func (*Progress) ProtoMessage()               {}
func (*Progress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Progress) GetQueryId() string {
	if m != nil {
		return m.QueryId
	}
	return ""
}

func (m *Progress) GetFilesProcessed() int64 {
	if m != nil {
		return m.FilesProcessed
	}
	return 0
}

func (m *Progress) GetFilesTotal() int64 {
	if m != nil {
		return m.FilesTotal
	}
	return 0
}

func (m *Progress) GetResults() int64 {
	if m != nil {
		return m.Results
	}
	return 0
}

type Pagination struct {
	QueryId     string `protobuf:"bytes,1,opt,name=query_id,json=queryId" json:"query_id,omitempty"`
	ResultPages int64  `protobuf:"varint,2,opt,name=result_pages,json=resultPages" json:"result_pages,omitempty"`
}

func (m *Pagination) Reset()                    { *m = Pagination{} }
func (m *Pagination) String() string            { return proto.CompactTextString(m) }
func (*Pagination) ProtoMessage()               {}
func (*Pagination) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Pagination) GetQueryId() string {
	if m != nil {
		return m.QueryId
	}
	return ""
}

func (m *Pagination) GetResultPages() int64 {
	if m != nil {
		return m.ResultPages
	}
	return 0
}

type Event struct {
	// Types that are valid to be assigned to Data:
	//	*Event_Error
	//	*Event_Progress
	//	*Event_Match
	//	*Event_Pagination
	Data isEvent_Data `protobuf_oneof:"data"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isEvent_Data interface {
	isEvent_Data()
}

type Event_Error struct {
	Error *Error `protobuf:"bytes,1,opt,name=error,oneof"`
}
type Event_Progress struct {
	Progress *Progress `protobuf:"bytes,2,opt,name=progress,oneof"`
}
type Event_Match struct {
	Match *sourcebackendpb.Match `protobuf:"bytes,3,opt,name=match,oneof"`
}
type Event_Pagination struct {
	Pagination *Pagination `protobuf:"bytes,4,opt,name=pagination,oneof"`
}

func (*Event_Error) isEvent_Data()      {}
func (*Event_Progress) isEvent_Data()   {}
func (*Event_Match) isEvent_Data()      {}
func (*Event_Pagination) isEvent_Data() {}

func (m *Event) GetData() isEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetError() *Error {
	if x, ok := m.GetData().(*Event_Error); ok {
		return x.Error
	}
	return nil
}

func (m *Event) GetProgress() *Progress {
	if x, ok := m.GetData().(*Event_Progress); ok {
		return x.Progress
	}
	return nil
}

func (m *Event) GetMatch() *sourcebackendpb.Match {
	if x, ok := m.GetData().(*Event_Match); ok {
		return x.Match
	}
	return nil
}

func (m *Event) GetPagination() *Pagination {
	if x, ok := m.GetData().(*Event_Pagination); ok {
		return x.Pagination
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Error)(nil),
		(*Event_Progress)(nil),
		(*Event_Match)(nil),
		(*Event_Pagination)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// data
	switch x := m.Data.(type) {
	case *Event_Error:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case *Event_Progress:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Progress); err != nil {
			return err
		}
	case *Event_Match:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Match); err != nil {
			return err
		}
	case *Event_Pagination:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Pagination); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Data has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 1: // data.error
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Error)
		err := b.DecodeMessage(msg)
		m.Data = &Event_Error{msg}
		return true, err
	case 2: // data.progress
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Progress)
		err := b.DecodeMessage(msg)
		m.Data = &Event_Progress{msg}
		return true, err
	case 3: // data.match
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(sourcebackendpb.Match)
		err := b.DecodeMessage(msg)
		m.Data = &Event_Match{msg}
		return true, err
	case 4: // data.pagination
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Pagination)
		err := b.DecodeMessage(msg)
		m.Data = &Event_Pagination{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// data
	switch x := m.Data.(type) {
	case *Event_Error:
		s := proto.Size(x.Error)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Progress:
		s := proto.Size(x.Progress)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Match:
		s := proto.Size(x.Match)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Pagination:
		s := proto.Size(x.Pagination)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*SearchRequest)(nil), "dcspb.SearchRequest")
	proto.RegisterType((*Error)(nil), "dcspb.Error")
	proto.RegisterType((*Progress)(nil), "dcspb.Progress")
	proto.RegisterType((*Pagination)(nil), "dcspb.Pagination")
	proto.RegisterType((*Event)(nil), "dcspb.Event")
	proto.RegisterEnum("dcspb.Error_ErrorType", Error_ErrorType_name, Error_ErrorType_value)
	proto.RegisterEnum("dcspb.Event_Type", Event_Type_name, Event_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DCS service

type DCSClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (DCS_SearchClient, error)
}

type dCSClient struct {
	cc *grpc.ClientConn
}

func NewDCSClient(cc *grpc.ClientConn) DCSClient {
	return &dCSClient{cc}
}

func (c *dCSClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (DCS_SearchClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DCS_serviceDesc.Streams[0], c.cc, "/dcspb.DCS/Search", opts...)
	if err != nil {
		return nil, err
	}
	x := &dCSSearchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DCS_SearchClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type dCSSearchClient struct {
	grpc.ClientStream
}

func (x *dCSSearchClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DCS service

type DCSServer interface {
	Search(*SearchRequest, DCS_SearchServer) error
}

func RegisterDCSServer(s *grpc.Server, srv DCSServer) {
	s.RegisterService(&_DCS_serviceDesc, srv)
}

func _DCS_Search_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DCSServer).Search(m, &dCSSearchServer{stream})
}

type DCS_SearchServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type dCSSearchServer struct {
	grpc.ServerStream
}

func (x *dCSSearchServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

var _DCS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dcspb.DCS",
	HandlerType: (*DCSServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Search",
			Handler:       _DCS_Search_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "dcs.proto",
}

func init() { proto.RegisterFile("dcs.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 517 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xb5, 0x13, 0x3b, 0x4d, 0xc6, 0x49, 0xea, 0x0e, 0x55, 0x09, 0xbd, 0x00, 0x06, 0x04, 0x42,
	0xc2, 0x54, 0xee, 0x81, 0xb3, 0x13, 0x9b, 0xc6, 0x90, 0x3a, 0x66, 0x93, 0x56, 0xe2, 0x14, 0x39,
	0xf6, 0x92, 0x46, 0xa4, 0xb1, 0xbb, 0xbb, 0x41, 0xca, 0x27, 0x70, 0xe2, 0x1b, 0xf8, 0x53, 0xe4,
	0x75, 0x62, 0xb5, 0x1c, 0x7a, 0xb1, 0x34, 0xef, 0x3d, 0xcf, 0xcc, 0x7b, 0x3b, 0xd0, 0x4a, 0x13,
	0x6e, 0xe7, 0x2c, 0x13, 0x19, 0xea, 0x69, 0xc2, 0xf3, 0xf9, 0xe9, 0x2b, 0x9e, 0x6d, 0x58, 0x42,
	0xe7, 0x71, 0xf2, 0x93, 0xae, 0xd3, 0x7c, 0xfe, 0xf1, 0x41, 0x5d, 0x6a, 0xad, 0x37, 0xd0, 0x99,
	0xd0, 0x98, 0x25, 0x37, 0x84, 0xde, 0x6d, 0x28, 0x17, 0x78, 0x0c, 0xfa, 0xdd, 0x86, 0xb2, 0x6d,
	0x4f, 0x7d, 0xa1, 0xbe, 0x6b, 0x91, 0xb2, 0xb0, 0xfe, 0xaa, 0xa0, 0xfb, 0x8c, 0x65, 0x0c, 0xdf,
	0x83, 0x26, 0xb6, 0x39, 0x95, 0x74, 0xd7, 0x39, 0xb1, 0xe5, 0x2c, 0x5b, 0x72, 0xe5, 0x77, 0xba,
	0xcd, 0x29, 0x91, 0x1a, 0xec, 0xc1, 0xc1, 0x2d, 0xe5, 0x3c, 0x5e, 0xd0, 0x5e, 0x4d, 0x76, 0xdb,
	0x97, 0x16, 0x81, 0x56, 0x25, 0xc6, 0x0e, 0xb4, 0x06, 0x6e, 0x38, 0xf0, 0x47, 0x23, 0xdf, 0x33,
	0x15, 0x7c, 0x0a, 0x4f, 0xfa, 0xee, 0xe0, 0xab, 0x1f, 0x7a, 0xb3, 0xab, 0xd0, 0xbd, 0x76, 0x83,
	0x91, 0xdb, 0x1f, 0xf9, 0xa6, 0x8a, 0x00, 0x8d, 0xcf, 0x6e, 0x50, 0x88, 0x6a, 0x78, 0x04, 0x9d,
	0x20, 0xbc, 0x76, 0x47, 0x81, 0x37, 0xfb, 0x76, 0xe5, 0x93, 0xef, 0x66, 0xdd, 0xfa, 0xad, 0x42,
	0x33, 0x62, 0xd9, 0x82, 0x51, 0xce, 0xf1, 0x19, 0x34, 0xe5, 0xe6, 0xb3, 0x65, 0xba, 0x73, 0x72,
	0x20, 0xeb, 0x20, 0xc5, 0xb7, 0x70, 0xf8, 0x63, 0xb9, 0xa2, 0x7c, 0x96, 0xb3, 0x2c, 0xa1, 0x9c,
	0xd3, 0x54, 0x6e, 0x57, 0x27, 0x5d, 0x09, 0x47, 0x7b, 0x14, 0x9f, 0x83, 0x51, 0x0a, 0x45, 0x26,
	0xe2, 0x55, 0xaf, 0x2e, 0x45, 0x20, 0xa1, 0x69, 0x81, 0x14, 0xfe, 0x18, 0xe5, 0x9b, 0x95, 0xe0,
	0x3d, 0x4d, 0x92, 0xfb, 0xd2, 0xfa, 0x02, 0x10, 0xc5, 0x8b, 0xe5, 0x3a, 0x16, 0xcb, 0x6c, 0xfd,
	0xd8, 0x32, 0x2f, 0xa1, 0x5d, 0xfe, 0x33, 0xcb, 0xe3, 0x05, 0xe5, 0xbb, 0x4d, 0x8c, 0x12, 0x8b,
	0x0a, 0xc8, 0xfa, 0x53, 0x03, 0xdd, 0xff, 0x45, 0xd7, 0x02, 0x5f, 0x83, 0x4e, 0x8b, 0xd4, 0x64,
	0x13, 0xc3, 0x69, 0xdf, 0x0f, 0x7f, 0xa8, 0x90, 0x92, 0xc4, 0x0f, 0xd0, 0xcc, 0x77, 0x31, 0xc8,
	0x76, 0x86, 0x73, 0xb8, 0x13, 0xee, 0xd3, 0x19, 0x2a, 0xa4, 0x92, 0xa0, 0x0d, 0xfa, 0x6d, 0x2c,
	0x92, 0x1b, 0xe9, 0xcf, 0x70, 0x4e, 0xec, 0xff, 0xce, 0xc6, 0xbe, 0x2c, 0xd8, 0xa2, 0xbd, 0x94,
	0xe1, 0x39, 0x40, 0x5e, 0x59, 0x93, 0xbe, 0x0d, 0xe7, 0x68, 0x3f, 0xa0, 0x22, 0x86, 0x0a, 0xb9,
	0x27, 0xb3, 0x3c, 0xd0, 0xe4, 0x53, 0xb7, 0x40, 0xf7, 0x09, 0x19, 0x13, 0x53, 0xc1, 0x36, 0x34,
	0x23, 0x32, 0xbe, 0x20, 0xfe, 0x64, 0x62, 0xaa, 0x05, 0x71, 0xe9, 0x4e, 0x07, 0x43, 0xb3, 0x86,
	0x5d, 0x80, 0xc8, 0xbd, 0x08, 0x42, 0x77, 0x1a, 0x8c, 0x43, 0xb3, 0x8e, 0x4d, 0xd0, 0xbc, 0x71,
	0xe8, 0x9b, 0x5a, 0xbf, 0x01, 0x5a, 0x1a, 0x8b, 0xd8, 0xf9, 0x04, 0x75, 0x6f, 0x30, 0xc1, 0x33,
	0x68, 0x94, 0xb7, 0x8b, 0xc7, 0xbb, 0xf9, 0x0f, 0x4e, 0xf9, 0xb4, 0xca, 0xa7, 0x08, 0xcf, 0x52,
	0xce, 0xd4, 0x79, 0x43, 0x1e, 0xfd, 0xf9, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6a, 0xde, 0xbd,
	0x38, 0x2d, 0x03, 0x00, 0x00,
}
