// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: app/router/command/command.proto

package command

import (
	proto "github.com/golang/protobuf/proto"
	net "github.com/whaleblueio/Xray-core/common/net"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// RoutingContext is the context with information relative to routing process.
// It conforms to the structure of xray.features.routing.Context and
// xray.features.routing.Route.
type RoutingContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InboundTag        string            `protobuf:"bytes,1,opt,name=InboundTag,proto3" json:"InboundTag,omitempty"`
	Network           net.Network       `protobuf:"varint,2,opt,name=Network,proto3,enum=xray.common.net.Network" json:"Network,omitempty"`
	SourceIPs         [][]byte          `protobuf:"bytes,3,rep,name=SourceIPs,proto3" json:"SourceIPs,omitempty"`
	TargetIPs         [][]byte          `protobuf:"bytes,4,rep,name=TargetIPs,proto3" json:"TargetIPs,omitempty"`
	SourcePort        uint32            `protobuf:"varint,5,opt,name=SourcePort,proto3" json:"SourcePort,omitempty"`
	TargetPort        uint32            `protobuf:"varint,6,opt,name=TargetPort,proto3" json:"TargetPort,omitempty"`
	TargetDomain      string            `protobuf:"bytes,7,opt,name=TargetDomain,proto3" json:"TargetDomain,omitempty"`
	Protocol          string            `protobuf:"bytes,8,opt,name=Protocol,proto3" json:"Protocol,omitempty"`
	User              string            `protobuf:"bytes,9,opt,name=User,proto3" json:"User,omitempty"`
	Attributes        map[string]string `protobuf:"bytes,10,rep,name=Attributes,proto3" json:"Attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	OutboundGroupTags []string          `protobuf:"bytes,11,rep,name=OutboundGroupTags,proto3" json:"OutboundGroupTags,omitempty"`
	OutboundTag       string            `protobuf:"bytes,12,opt,name=OutboundTag,proto3" json:"OutboundTag,omitempty"`
}

func (x *RoutingContext) Reset() {
	*x = RoutingContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_router_command_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoutingContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoutingContext) ProtoMessage() {}

func (x *RoutingContext) ProtoReflect() protoreflect.Message {
	mi := &file_app_router_command_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoutingContext.ProtoReflect.Descriptor instead.
func (*RoutingContext) Descriptor() ([]byte, []int) {
	return file_app_router_command_command_proto_rawDescGZIP(), []int{0}
}

func (x *RoutingContext) GetInboundTag() string {
	if x != nil {
		return x.InboundTag
	}
	return ""
}

func (x *RoutingContext) GetNetwork() net.Network {
	if x != nil {
		return x.Network
	}
	return net.Network_Unknown
}

func (x *RoutingContext) GetSourceIPs() [][]byte {
	if x != nil {
		return x.SourceIPs
	}
	return nil
}

func (x *RoutingContext) GetTargetIPs() [][]byte {
	if x != nil {
		return x.TargetIPs
	}
	return nil
}

func (x *RoutingContext) GetSourcePort() uint32 {
	if x != nil {
		return x.SourcePort
	}
	return 0
}

func (x *RoutingContext) GetTargetPort() uint32 {
	if x != nil {
		return x.TargetPort
	}
	return 0
}

func (x *RoutingContext) GetTargetDomain() string {
	if x != nil {
		return x.TargetDomain
	}
	return ""
}

func (x *RoutingContext) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *RoutingContext) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *RoutingContext) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *RoutingContext) GetOutboundGroupTags() []string {
	if x != nil {
		return x.OutboundGroupTags
	}
	return nil
}

func (x *RoutingContext) GetOutboundTag() string {
	if x != nil {
		return x.OutboundTag
	}
	return ""
}

// SubscribeRoutingStatsRequest subscribes to routing statistics channel if
// opened by xray-core.
// * FieldSelectors selects a subset of fields in routing statistics to return.
// Valid selectors:
//  - inbound: Selects connection's inbound tag.
//  - network: Selects connection's network.
//  - ip: Equivalent as "ip_source" and "ip_target", selects both source and
//  target IP.
//  - port: Equivalent as "port_source" and "port_target", selects both source
//  and target port.
//  - domain: Selects target domain.
//  - protocol: Select connection's protocol.
//  - user: Select connection's inbound user email.
//  - attributes: Select connection's additional attributes.
//  - outbound: Equivalent as "outbound" and "outbound_group", select both
//  outbound tag and outbound group tags.
// * If FieldSelectors is left empty, all fields will be returned.
type SubscribeRoutingStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldSelectors []string `protobuf:"bytes,1,rep,name=FieldSelectors,proto3" json:"FieldSelectors,omitempty"`
}

func (x *SubscribeRoutingStatsRequest) Reset() {
	*x = SubscribeRoutingStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_router_command_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRoutingStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRoutingStatsRequest) ProtoMessage() {}

func (x *SubscribeRoutingStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_router_command_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRoutingStatsRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRoutingStatsRequest) Descriptor() ([]byte, []int) {
	return file_app_router_command_command_proto_rawDescGZIP(), []int{1}
}

func (x *SubscribeRoutingStatsRequest) GetFieldSelectors() []string {
	if x != nil {
		return x.FieldSelectors
	}
	return nil
}

// TestRouteRequest manually tests a routing result according to the routing
// context message.
// * RoutingContext is the routing message without outbound information.
// * FieldSelectors selects the fields to return in the routing result. All
// fields are returned if left empty.
// * PublishResult broadcasts the routing result to routing statistics channel
// if set true.
type TestRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoutingContext *RoutingContext `protobuf:"bytes,1,opt,name=RoutingContext,proto3" json:"RoutingContext,omitempty"`
	FieldSelectors []string        `protobuf:"bytes,2,rep,name=FieldSelectors,proto3" json:"FieldSelectors,omitempty"`
	PublishResult  bool            `protobuf:"varint,3,opt,name=PublishResult,proto3" json:"PublishResult,omitempty"`
}

func (x *TestRouteRequest) Reset() {
	*x = TestRouteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_router_command_command_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRouteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRouteRequest) ProtoMessage() {}

func (x *TestRouteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_router_command_command_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRouteRequest.ProtoReflect.Descriptor instead.
func (*TestRouteRequest) Descriptor() ([]byte, []int) {
	return file_app_router_command_command_proto_rawDescGZIP(), []int{2}
}

func (x *TestRouteRequest) GetRoutingContext() *RoutingContext {
	if x != nil {
		return x.RoutingContext
	}
	return nil
}

func (x *TestRouteRequest) GetFieldSelectors() []string {
	if x != nil {
		return x.FieldSelectors
	}
	return nil
}

func (x *TestRouteRequest) GetPublishResult() bool {
	if x != nil {
		return x.PublishResult
	}
	return false
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_router_command_command_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_app_router_command_command_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_app_router_command_command_proto_rawDescGZIP(), []int{3}
}

var File_app_router_command_command_proto protoreflect.FileDescriptor

var file_app_router_command_command_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x70, 0x70, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x17, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x18, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x04, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x49, 0x6e, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x54, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x49, 0x6e,
	0x62, 0x6f, 0x75, 0x6e, 0x64, 0x54, 0x61, 0x67, 0x12, 0x32, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x78, 0x72, 0x61, 0x79,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x52, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x50, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52,
	0x09, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x50, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x49, 0x50, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x09, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x50, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x57, 0x0a, 0x0a,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x37, 0x2e, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e,
	0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x61, 0x67, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x11, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54,
	0x61, 0x67, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x54,
	0x61, 0x67, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75,
	0x6e, 0x64, 0x54, 0x61, 0x67, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x46, 0x0a, 0x1c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x22, 0xb1, 0x01, 0x0a,
	0x10, 0x54, 0x65, 0x73, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x4f, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x78, 0x72, 0x61, 0x79,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x52, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0d, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x08, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x32, 0xf0, 0x01, 0x0a, 0x0e, 0x52,
	0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7b, 0x0a,
	0x15, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x35, 0x2e, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e,
	0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12, 0x61, 0x0a, 0x09, 0x54, 0x65,
	0x73, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x29, 0x2e, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x27, 0x2e, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x22, 0x00, 0x42, 0x67, 0x0a,
	0x1b, 0x63, 0x6f, 0x6d, 0x2e, 0x78, 0x72, 0x61, 0x79, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x50, 0x01, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x74, 0x6c, 0x73, 0x2f,
	0x78, 0x72, 0x61, 0x79, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0xaa, 0x02, 0x17, 0x58,
	0x72, 0x61, 0x79, 0x2e, 0x41, 0x70, 0x70, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_router_command_command_proto_rawDescOnce sync.Once
	file_app_router_command_command_proto_rawDescData = file_app_router_command_command_proto_rawDesc
)

func file_app_router_command_command_proto_rawDescGZIP() []byte {
	file_app_router_command_command_proto_rawDescOnce.Do(func() {
		file_app_router_command_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_router_command_command_proto_rawDescData)
	})
	return file_app_router_command_command_proto_rawDescData
}

var file_app_router_command_command_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_app_router_command_command_proto_goTypes = []interface{}{
	(*RoutingContext)(nil),               // 0: xray.app.router.command.RoutingContext
	(*SubscribeRoutingStatsRequest)(nil), // 1: xray.app.router.command.SubscribeRoutingStatsRequest
	(*TestRouteRequest)(nil),             // 2: xray.app.router.command.TestRouteRequest
	(*Config)(nil),                       // 3: xray.app.router.command.Config
	nil,                                  // 4: xray.app.router.command.RoutingContext.AttributesEntry
	(net.Network)(0),                     // 5: xray.common.net.Network
}
var file_app_router_command_command_proto_depIdxs = []int32{
	5, // 0: xray.app.router.command.RoutingContext.Network:type_name -> xray.common.net.Network
	4, // 1: xray.app.router.command.RoutingContext.Attributes:type_name -> xray.app.router.command.RoutingContext.AttributesEntry
	0, // 2: xray.app.router.command.TestRouteRequest.RoutingContext:type_name -> xray.app.router.command.RoutingContext
	1, // 3: xray.app.router.command.RoutingService.SubscribeRoutingStats:input_type -> xray.app.router.command.SubscribeRoutingStatsRequest
	2, // 4: xray.app.router.command.RoutingService.TestRoute:input_type -> xray.app.router.command.TestRouteRequest
	0, // 5: xray.app.router.command.RoutingService.SubscribeRoutingStats:output_type -> xray.app.router.command.RoutingContext
	0, // 6: xray.app.router.command.RoutingService.TestRoute:output_type -> xray.app.router.command.RoutingContext
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_app_router_command_command_proto_init() }
func file_app_router_command_command_proto_init() {
	if File_app_router_command_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_router_command_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoutingContext); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_router_command_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeRoutingStatsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_router_command_command_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRouteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_router_command_command_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_router_command_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_app_router_command_command_proto_goTypes,
		DependencyIndexes: file_app_router_command_command_proto_depIdxs,
		MessageInfos:      file_app_router_command_command_proto_msgTypes,
	}.Build()
	File_app_router_command_command_proto = out.File
	file_app_router_command_command_proto_rawDesc = nil
	file_app_router_command_command_proto_goTypes = nil
	file_app_router_command_command_proto_depIdxs = nil
}
