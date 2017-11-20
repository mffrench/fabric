// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer/events.proto

package peer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/hyperledger/fabric/protos/common"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EventType int32

const (
	EventType_REGISTER      EventType = 0
	EventType_BLOCK         EventType = 1
	EventType_CHAINCODE     EventType = 2
	EventType_REJECTION     EventType = 3
	EventType_FILTEREDBLOCK EventType = 4
)

var EventType_name = map[int32]string{
	0: "REGISTER",
	1: "BLOCK",
	2: "CHAINCODE",
	3: "REJECTION",
	4: "FILTEREDBLOCK",
}
var EventType_value = map[string]int32{
	"REGISTER":      0,
	"BLOCK":         1,
	"CHAINCODE":     2,
	"REJECTION":     3,
	"FILTEREDBLOCK": 4,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

// ChaincodeReg is used for registering chaincode Interests
// when EventType is CHAINCODE
type ChaincodeReg struct {
	ChaincodeId string `protobuf:"bytes,1,opt,name=chaincode_id,json=chaincodeId" json:"chaincode_id,omitempty"`
	EventName   string `protobuf:"bytes,2,opt,name=event_name,json=eventName" json:"event_name,omitempty"`
}

func (m *ChaincodeReg) Reset()                    { *m = ChaincodeReg{} }
func (m *ChaincodeReg) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeReg) ProtoMessage()               {}
func (*ChaincodeReg) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

func (m *ChaincodeReg) GetChaincodeId() string {
	if m != nil {
		return m.ChaincodeId
	}
	return ""
}

func (m *ChaincodeReg) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

type Interest struct {
	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,enum=protos.EventType" json:"event_type,omitempty"`
	// Ideally we should just have the following oneof for different
	// Reg types and get rid of EventType. But this is an API change
	// Additional Reg types may add messages specific to their type
	// to the oneof.
	//
	// Types that are valid to be assigned to RegInfo:
	//	*Interest_ChaincodeRegInfo
	RegInfo isInterest_RegInfo `protobuf_oneof:"RegInfo"`
	ChainID string             `protobuf:"bytes,3,opt,name=chainID" json:"chainID,omitempty"`
}

func (m *Interest) Reset()                    { *m = Interest{} }
func (m *Interest) String() string            { return proto.CompactTextString(m) }
func (*Interest) ProtoMessage()               {}
func (*Interest) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{1} }

type isInterest_RegInfo interface {
	isInterest_RegInfo()
}

type Interest_ChaincodeRegInfo struct {
	ChaincodeRegInfo *ChaincodeReg `protobuf:"bytes,2,opt,name=chaincode_reg_info,json=chaincodeRegInfo,oneof"`
}

func (*Interest_ChaincodeRegInfo) isInterest_RegInfo() {}

func (m *Interest) GetRegInfo() isInterest_RegInfo {
	if m != nil {
		return m.RegInfo
	}
	return nil
}

func (m *Interest) GetEventType() EventType {
	if m != nil {
		return m.EventType
	}
	return EventType_REGISTER
}

func (m *Interest) GetChaincodeRegInfo() *ChaincodeReg {
	if x, ok := m.GetRegInfo().(*Interest_ChaincodeRegInfo); ok {
		return x.ChaincodeRegInfo
	}
	return nil
}

func (m *Interest) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Interest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Interest_OneofMarshaler, _Interest_OneofUnmarshaler, _Interest_OneofSizer, []interface{}{
		(*Interest_ChaincodeRegInfo)(nil),
	}
}

func _Interest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeRegInfo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Interest.RegInfo has unexpected type %T", x)
	}
	return nil
}

func _Interest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Interest)
	switch tag {
	case 2: // RegInfo.chaincode_reg_info
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeReg)
		err := b.DecodeMessage(msg)
		m.RegInfo = &Interest_ChaincodeRegInfo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Interest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		s := proto.Size(x.ChaincodeRegInfo)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ---------- consumer events ---------
// Register is sent by consumers for registering events
// string type - "register"
type Register struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Register) Reset()                    { *m = Register{} }
func (m *Register) String() string            { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()               {}
func (*Register) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2} }

func (m *Register) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// Rejection is sent by consumers for erroneous transaction rejection events
// string type - "rejection"
type Rejection struct {
	Tx       *Transaction `protobuf:"bytes,1,opt,name=tx" json:"tx,omitempty"`
	ErrorMsg string       `protobuf:"bytes,2,opt,name=error_msg,json=errorMsg" json:"error_msg,omitempty"`
}

func (m *Rejection) Reset()                    { *m = Rejection{} }
func (m *Rejection) String() string            { return proto.CompactTextString(m) }
func (*Rejection) ProtoMessage()               {}
func (*Rejection) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{3} }

func (m *Rejection) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *Rejection) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

// ---------- producer events ---------
type Unregister struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Unregister) Reset()                    { *m = Unregister{} }
func (m *Unregister) String() string            { return proto.CompactTextString(m) }
func (*Unregister) ProtoMessage()               {}
func (*Unregister) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{4} }

func (m *Unregister) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// FilteredBlock is sent by producers and contains minimal information
// about the block.
type FilteredBlock struct {
	ChannelId  string                 `protobuf:"bytes,1,opt,name=channel_id,json=channelId" json:"channel_id,omitempty"`
	Number     uint64                 `protobuf:"varint,2,opt,name=number" json:"number,omitempty"`
	Type       common.HeaderType      `protobuf:"varint,3,opt,name=type,enum=common.HeaderType" json:"type,omitempty"`
	FilteredTx []*FilteredTransaction `protobuf:"bytes,4,rep,name=filtered_tx,json=filteredTx" json:"filtered_tx,omitempty"`
}

func (m *FilteredBlock) Reset()                    { *m = FilteredBlock{} }
func (m *FilteredBlock) String() string            { return proto.CompactTextString(m) }
func (*FilteredBlock) ProtoMessage()               {}
func (*FilteredBlock) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{5} }

func (m *FilteredBlock) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *FilteredBlock) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *FilteredBlock) GetType() common.HeaderType {
	if m != nil {
		return m.Type
	}
	return common.HeaderType_MESSAGE
}

func (m *FilteredBlock) GetFilteredTx() []*FilteredTransaction {
	if m != nil {
		return m.FilteredTx
	}
	return nil
}

// FilteredTransaction is a minimal set of information about a transaction
// within a block.
type FilteredTransaction struct {
	Txid             string            `protobuf:"bytes,1,opt,name=txid" json:"txid,omitempty"`
	TxValidationCode TxValidationCode  `protobuf:"varint,2,opt,name=tx_validation_code,json=txValidationCode,enum=protos.TxValidationCode" json:"tx_validation_code,omitempty"`
	FilteredAction   []*FilteredAction `protobuf:"bytes,3,rep,name=filtered_action,json=filteredAction" json:"filtered_action,omitempty"`
}

func (m *FilteredTransaction) Reset()                    { *m = FilteredTransaction{} }
func (m *FilteredTransaction) String() string            { return proto.CompactTextString(m) }
func (*FilteredTransaction) ProtoMessage()               {}
func (*FilteredTransaction) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{6} }

func (m *FilteredTransaction) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func (m *FilteredTransaction) GetTxValidationCode() TxValidationCode {
	if m != nil {
		return m.TxValidationCode
	}
	return TxValidationCode_VALID
}

func (m *FilteredTransaction) GetFilteredAction() []*FilteredAction {
	if m != nil {
		return m.FilteredAction
	}
	return nil
}

// FilteredAction is a minimal set of information about an action within a
// transaction.
type FilteredAction struct {
	CcEvent *ChaincodeEvent `protobuf:"bytes,1,opt,name=ccEvent" json:"ccEvent,omitempty"`
}

func (m *FilteredAction) Reset()                    { *m = FilteredAction{} }
func (m *FilteredAction) String() string            { return proto.CompactTextString(m) }
func (*FilteredAction) ProtoMessage()               {}
func (*FilteredAction) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{7} }

func (m *FilteredAction) GetCcEvent() *ChaincodeEvent {
	if m != nil {
		return m.CcEvent
	}
	return nil
}

// SignedEvent is used for any communication between consumer and producer
type SignedEvent struct {
	// Signature over the event bytes
	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// Marshal of Event object
	EventBytes []byte `protobuf:"bytes,2,opt,name=eventBytes,proto3" json:"eventBytes,omitempty"`
}

func (m *SignedEvent) Reset()                    { *m = SignedEvent{} }
func (m *SignedEvent) String() string            { return proto.CompactTextString(m) }
func (*SignedEvent) ProtoMessage()               {}
func (*SignedEvent) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{8} }

func (m *SignedEvent) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SignedEvent) GetEventBytes() []byte {
	if m != nil {
		return m.EventBytes
	}
	return nil
}

// Event is used by
//  - consumers (adapters) to send Register
//  - producer to advertise supported types and events
type Event struct {
	// Types that are valid to be assigned to Event:
	//	*Event_Register
	//	*Event_Block
	//	*Event_ChaincodeEvent
	//	*Event_Rejection
	//	*Event_Unregister
	//	*Event_FilteredBlock
	Event isEvent_Event `protobuf_oneof:"Event"`
	// Creator of the event, specified as a certificate chain
	Creator []byte `protobuf:"bytes,6,opt,name=creator,proto3" json:"creator,omitempty"`
	// Timestamp of the client - used to mitigate replay attacks
	Timestamp *google_protobuf1.Timestamp `protobuf:"bytes,8,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{9} }

type isEvent_Event interface {
	isEvent_Event()
}

type Event_Register struct {
	Register *Register `protobuf:"bytes,1,opt,name=register,oneof"`
}
type Event_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,oneof"`
}
type Event_ChaincodeEvent struct {
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,3,opt,name=chaincode_event,json=chaincodeEvent,oneof"`
}
type Event_Rejection struct {
	Rejection *Rejection `protobuf:"bytes,4,opt,name=rejection,oneof"`
}
type Event_Unregister struct {
	Unregister *Unregister `protobuf:"bytes,5,opt,name=unregister,oneof"`
}
type Event_FilteredBlock struct {
	FilteredBlock *FilteredBlock `protobuf:"bytes,7,opt,name=filtered_block,json=filteredBlock,oneof"`
}

func (*Event_Register) isEvent_Event()       {}
func (*Event_Block) isEvent_Event()          {}
func (*Event_ChaincodeEvent) isEvent_Event() {}
func (*Event_Rejection) isEvent_Event()      {}
func (*Event_Unregister) isEvent_Event()     {}
func (*Event_FilteredBlock) isEvent_Event()  {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetRegister() *Register {
	if x, ok := m.GetEvent().(*Event_Register); ok {
		return x.Register
	}
	return nil
}

func (m *Event) GetBlock() *common.Block {
	if x, ok := m.GetEvent().(*Event_Block); ok {
		return x.Block
	}
	return nil
}

func (m *Event) GetChaincodeEvent() *ChaincodeEvent {
	if x, ok := m.GetEvent().(*Event_ChaincodeEvent); ok {
		return x.ChaincodeEvent
	}
	return nil
}

func (m *Event) GetRejection() *Rejection {
	if x, ok := m.GetEvent().(*Event_Rejection); ok {
		return x.Rejection
	}
	return nil
}

func (m *Event) GetUnregister() *Unregister {
	if x, ok := m.GetEvent().(*Event_Unregister); ok {
		return x.Unregister
	}
	return nil
}

func (m *Event) GetFilteredBlock() *FilteredBlock {
	if x, ok := m.GetEvent().(*Event_FilteredBlock); ok {
		return x.FilteredBlock
	}
	return nil
}

func (m *Event) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *Event) GetTimestamp() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Register)(nil),
		(*Event_Block)(nil),
		(*Event_ChaincodeEvent)(nil),
		(*Event_Rejection)(nil),
		(*Event_Unregister)(nil),
		(*Event_FilteredBlock)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Register); err != nil {
			return err
		}
	case *Event_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *Event_ChaincodeEvent:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeEvent); err != nil {
			return err
		}
	case *Event_Rejection:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Rejection); err != nil {
			return err
		}
	case *Event_Unregister:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unregister); err != nil {
			return err
		}
	case *Event_FilteredBlock:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.FilteredBlock); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 1: // Event.register
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Register)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Register{msg}
		return true, err
	case 2: // Event.block
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.Block)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Block{msg}
		return true, err
	case 3: // Event.chaincode_event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeEvent)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ChaincodeEvent{msg}
		return true, err
	case 4: // Event.rejection
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Rejection)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Rejection{msg}
		return true, err
	case 5: // Event.unregister
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Unregister)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Unregister{msg}
		return true, err
	case 7: // Event.filtered_block
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FilteredBlock)
		err := b.DecodeMessage(msg)
		m.Event = &Event_FilteredBlock{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		s := proto.Size(x.Register)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Block:
		s := proto.Size(x.Block)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ChaincodeEvent:
		s := proto.Size(x.ChaincodeEvent)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Rejection:
		s := proto.Size(x.Rejection)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Unregister:
		s := proto.Size(x.Unregister)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_FilteredBlock:
		s := proto.Size(x.FilteredBlock)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ChaincodeReg)(nil), "protos.ChaincodeReg")
	proto.RegisterType((*Interest)(nil), "protos.Interest")
	proto.RegisterType((*Register)(nil), "protos.Register")
	proto.RegisterType((*Rejection)(nil), "protos.Rejection")
	proto.RegisterType((*Unregister)(nil), "protos.Unregister")
	proto.RegisterType((*FilteredBlock)(nil), "protos.FilteredBlock")
	proto.RegisterType((*FilteredTransaction)(nil), "protos.FilteredTransaction")
	proto.RegisterType((*FilteredAction)(nil), "protos.FilteredAction")
	proto.RegisterType((*SignedEvent)(nil), "protos.SignedEvent")
	proto.RegisterType((*Event)(nil), "protos.Event")
	proto.RegisterEnum("protos.EventType", EventType_name, EventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Events service

type EventsClient interface {
	// event chatting using Event
	Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error)
}

type eventsClient struct {
	cc *grpc.ClientConn
}

func NewEventsClient(cc *grpc.ClientConn) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Events_serviceDesc.Streams[0], c.cc, "/protos.Events/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsChatClient{stream}
	return x, nil
}

type Events_ChatClient interface {
	Send(*SignedEvent) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsChatClient struct {
	grpc.ClientStream
}

func (x *eventsChatClient) Send(m *SignedEvent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventsChatClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Events service

type EventsServer interface {
	// event chatting using Event
	Chat(Events_ChatServer) error
}

func RegisterEventsServer(s *grpc.Server, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventsServer).Chat(&eventsChatServer{stream})
}

type Events_ChatServer interface {
	Send(*Event) error
	Recv() (*SignedEvent, error)
	grpc.ServerStream
}

type eventsChatServer struct {
	grpc.ServerStream
}

func (x *eventsChatServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventsChatServer) Recv() (*SignedEvent, error) {
	m := new(SignedEvent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _Events_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/events.proto",
}

func init() { proto.RegisterFile("peer/events.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 859 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdd, 0x8f, 0xdb, 0x44,
	0x10, 0x8f, 0x2f, 0xb9, 0x24, 0x9e, 0x7c, 0x34, 0x37, 0x07, 0x27, 0x2b, 0xe5, 0xa3, 0x18, 0x81,
	0x0e, 0x1e, 0x9c, 0x23, 0x54, 0x08, 0x21, 0x04, 0xba, 0xe4, 0x72, 0xc4, 0xb4, 0xbd, 0xab, 0xb6,
	0x29, 0x0f, 0x3c, 0x10, 0x39, 0xf6, 0xc6, 0x31, 0x8d, 0xed, 0x68, 0xbd, 0xa9, 0x72, 0x7f, 0x51,
	0x5f, 0x78, 0xe2, 0x2f, 0x44, 0x1e, 0x7b, 0xed, 0x24, 0xc0, 0x43, 0x9f, 0xec, 0x99, 0xf9, 0xcd,
	0xec, 0x6f, 0x76, 0x3e, 0x16, 0xce, 0x36, 0x9c, 0x8b, 0x01, 0x7f, 0xcb, 0x23, 0x99, 0x58, 0x1b,
	0x11, 0xcb, 0x18, 0xeb, 0xf4, 0x49, 0xfa, 0xe7, 0x6e, 0x1c, 0x86, 0x71, 0x34, 0xc8, 0x3e, 0x99,
	0xb1, 0xff, 0xa9, 0x1f, 0xc7, 0xfe, 0x9a, 0x0f, 0x48, 0x5a, 0x6c, 0x97, 0x03, 0x19, 0x84, 0x3c,
	0x91, 0x4e, 0xb8, 0xc9, 0x01, 0x7d, 0x0a, 0xe8, 0xae, 0x9c, 0x20, 0x72, 0x63, 0x8f, 0xcf, 0x29,
	0x74, 0x6e, 0xbb, 0x20, 0x9b, 0x14, 0x4e, 0x94, 0x38, 0xae, 0x0c, 0x54, 0x50, 0xf3, 0x25, 0xb4,
	0xc7, 0xca, 0x81, 0x71, 0x1f, 0x3f, 0x83, 0x76, 0x19, 0x20, 0xf0, 0x0c, 0xed, 0x89, 0x76, 0xa9,
	0xb3, 0x56, 0xa1, 0xb3, 0x3d, 0xfc, 0x18, 0x80, 0x22, 0xcf, 0x23, 0x27, 0xe4, 0xc6, 0x09, 0x01,
	0x74, 0xd2, 0xdc, 0x39, 0x21, 0x37, 0xdf, 0x69, 0xd0, 0xb4, 0x23, 0xc9, 0x05, 0x4f, 0x24, 0x5e,
	0x29, 0xac, 0x7c, 0xd8, 0x70, 0x0a, 0xd6, 0x1d, 0x9e, 0x65, 0x47, 0x27, 0xd6, 0x24, 0xb5, 0xcc,
	0x1e, 0x36, 0x3c, 0x77, 0x4f, 0x7f, 0xf1, 0x06, 0xb0, 0x24, 0x20, 0xb8, 0x3f, 0x0f, 0xa2, 0x65,
	0x4c, 0xa7, 0xb4, 0x86, 0x1f, 0x28, 0xcf, 0x7d, 0xca, 0xd3, 0x0a, 0xeb, 0xb9, 0x7b, 0xb2, 0x1d,
	0x2d, 0x63, 0x34, 0xa0, 0x41, 0x3a, 0xfb, 0xc6, 0xa8, 0x12, 0x41, 0x25, 0x8e, 0x74, 0x68, 0xe4,
	0x20, 0xf3, 0x29, 0x34, 0x19, 0xf7, 0x83, 0x44, 0x72, 0x81, 0x97, 0x50, 0xcf, 0x2a, 0x61, 0x68,
	0x4f, 0xaa, 0x97, 0xad, 0x61, 0x4f, 0x1d, 0xa5, 0x52, 0x61, 0xb9, 0xdd, 0x7c, 0x01, 0x3a, 0xe3,
	0x7f, 0x72, 0xba, 0x44, 0xfc, 0x1c, 0x4e, 0xe4, 0x8e, 0xf2, 0x6a, 0x0d, 0xcf, 0x95, 0xcb, 0xac,
	0xbc, 0x65, 0x76, 0x22, 0x77, 0xf8, 0x18, 0x74, 0x2e, 0x44, 0x2c, 0xe6, 0x61, 0xe2, 0xe7, 0xf7,
	0xd5, 0x24, 0xc5, 0x8b, 0xc4, 0x37, 0xbf, 0x03, 0x78, 0x1d, 0x89, 0xf7, 0xa7, 0xf1, 0x97, 0x06,
	0x9d, 0xdb, 0x60, 0x9d, 0x6a, 0xbd, 0xd1, 0x3a, 0x76, 0xdf, 0xa4, 0x75, 0x71, 0x57, 0x4e, 0x14,
	0xf1, 0x75, 0x59, 0x38, 0x3d, 0xd7, 0xd8, 0x1e, 0x5e, 0x40, 0x3d, 0xda, 0x86, 0x0b, 0x2e, 0x88,
	0x42, 0x8d, 0xe5, 0x12, 0x7e, 0x09, 0x35, 0x2a, 0x4e, 0x95, 0x8a, 0x83, 0x56, 0xde, 0x73, 0x53,
	0xee, 0x78, 0x5c, 0x50, 0x75, 0xc8, 0x8e, 0x3f, 0x42, 0x6b, 0x99, 0x9f, 0x37, 0x97, 0x3b, 0xa3,
	0x46, 0xfc, 0x1e, 0x2b, 0x7e, 0x8a, 0xca, 0x7e, 0xee, 0xa0, 0xf0, 0xb3, 0x9d, 0xf9, 0xb7, 0x06,
	0xe7, 0xff, 0x81, 0x41, 0x84, 0x9a, 0xdc, 0x15, 0x74, 0xe9, 0x1f, 0x6f, 0x01, 0xe5, 0x6e, 0xfe,
	0xd6, 0x59, 0x07, 0x9e, 0x93, 0x82, 0xe6, 0x69, 0x65, 0x89, 0x75, 0x77, 0x68, 0x14, 0x97, 0xbc,
	0xfb, 0xad, 0x00, 0x8c, 0xd3, 0xca, 0xf7, 0xe4, 0x91, 0x06, 0x7f, 0x86, 0x47, 0x05, 0xe3, 0xec,
	0x38, 0xa3, 0x4a, 0xac, 0x2f, 0x8e, 0x59, 0x5f, 0x67, 0x84, 0xbb, 0xcb, 0x03, 0xd9, 0x1c, 0x41,
	0xf7, 0x10, 0x81, 0x57, 0xd0, 0x70, 0x5d, 0xea, 0xdb, 0xbc, 0xe8, 0x17, 0xff, 0x6a, 0x49, 0xb2,
	0x32, 0x05, 0x33, 0x9f, 0x41, 0xeb, 0x55, 0xe0, 0x47, 0xdc, 0x23, 0x11, 0x3f, 0x02, 0x3d, 0x09,
	0xfc, 0xc8, 0x91, 0x5b, 0x91, 0xcd, 0x43, 0x9b, 0x95, 0x0a, 0xfc, 0x24, 0x1f, 0x97, 0xd1, 0x83,
	0xe4, 0x09, 0x65, 0xdc, 0x66, 0x7b, 0x1a, 0xf3, 0x5d, 0x15, 0x4e, 0xb3, 0x38, 0x16, 0x34, 0x55,
	0xd3, 0xe4, 0x4c, 0x8a, 0x56, 0x51, 0x3d, 0x3d, 0xad, 0xb0, 0x02, 0x83, 0x5f, 0xc0, 0xe9, 0x22,
	0xed, 0x92, 0x7c, 0x92, 0x3a, 0xaa, 0xcc, 0xd4, 0x3a, 0xd3, 0x0a, 0xcb, 0xac, 0x78, 0x0d, 0x8f,
	0x8e, 0xf6, 0x07, 0xf5, 0xc5, 0xff, 0xe6, 0x39, 0xad, 0xb0, 0xae, 0x7b, 0xa0, 0xc1, 0x6f, 0x40,
	0x17, 0x6a, 0x3e, 0x8c, 0x1a, 0x39, 0x9f, 0x95, 0xd4, 0x72, 0xc3, 0xb4, 0xc2, 0x4a, 0x14, 0x3e,
	0x05, 0xd8, 0x16, 0x33, 0x60, 0x9c, 0x92, 0x0f, 0x2a, 0x9f, 0x72, 0x3a, 0xa6, 0x15, 0xb6, 0x87,
	0xc3, 0x9f, 0xa0, 0xa8, 0xd7, 0x3c, 0xcb, 0xad, 0x41, 0x9e, 0x1f, 0x1e, 0x57, 0x57, 0xe5, 0xd8,
	0x59, 0x1e, 0xcc, 0x4b, 0xba, 0x23, 0x04, 0x77, 0x64, 0x2c, 0x8c, 0x3a, 0xdd, 0xb4, 0x12, 0xf1,
	0x7b, 0xd0, 0x8b, 0xdd, 0x6a, 0x34, 0x29, 0x68, 0xdf, 0xca, 0xb6, 0xaf, 0xa5, 0xb6, 0xaf, 0x35,
	0x53, 0x08, 0x56, 0x82, 0x47, 0x8d, 0xbc, 0x3e, 0x5f, 0xbf, 0x06, 0xbd, 0x58, 0x6f, 0xd8, 0x86,
	0x26, 0x9b, 0xfc, 0x62, 0xbf, 0x9a, 0x4d, 0x58, 0xaf, 0x82, 0x3a, 0x9c, 0x8e, 0x9e, 0xdf, 0x8f,
	0x9f, 0xf5, 0x34, 0xec, 0x80, 0x3e, 0x9e, 0x5e, 0xdb, 0x77, 0xe3, 0xfb, 0x9b, 0x49, 0xef, 0x24,
	0x15, 0xd9, 0xe4, 0xd7, 0xc9, 0x78, 0x66, 0xdf, 0xdf, 0xf5, 0xaa, 0x78, 0x06, 0x9d, 0x5b, 0xfb,
	0xf9, 0x6c, 0xc2, 0x26, 0x37, 0x99, 0x43, 0x6d, 0xf8, 0x03, 0xd4, 0x29, 0x6c, 0x82, 0x57, 0x50,
	0x1b, 0xaf, 0x1c, 0x89, 0xc5, 0xd6, 0xd9, 0xeb, 0xb2, 0x7e, 0xe7, 0x60, 0xc5, 0x9a, 0x95, 0x4b,
	0xed, 0x4a, 0x1b, 0xfd, 0x01, 0x66, 0x2c, 0x7c, 0x6b, 0xf5, 0xb0, 0xe1, 0x62, 0xcd, 0x3d, 0x9f,
	0x0b, 0x6b, 0xe9, 0x2c, 0x44, 0xe0, 0x2a, 0x70, 0xfa, 0x44, 0x8c, 0x3a, 0x59, 0xfc, 0x97, 0x8e,
	0xfb, 0xc6, 0xf1, 0xf9, 0xef, 0x5f, 0xf9, 0x81, 0x5c, 0x6d, 0x17, 0x69, 0xbb, 0x0c, 0xf6, 0x3c,
	0x07, 0x99, 0x67, 0xf6, 0x16, 0x25, 0x83, 0xd4, 0x73, 0x91, 0x3d, 0x5e, 0xdf, 0xfe, 0x13, 0x00,
	0x00, 0xff, 0xff, 0x19, 0x90, 0x5a, 0x10, 0xd8, 0x06, 0x00, 0x00,
}
