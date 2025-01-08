// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.1
// source: chat/v1/chat.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetOnlineUserIdsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerId    uint32                 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOnlineUserIdsRequest) Reset() {
	*x = GetOnlineUserIdsRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUserIdsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUserIdsRequest) ProtoMessage() {}

func (x *GetOnlineUserIdsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOnlineUserIdsRequest.ProtoReflect.Descriptor instead.
func (*GetOnlineUserIdsRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{0}
}

func (x *GetOnlineUserIdsRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *GetOnlineUserIdsRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type GetOnlineUserIdsReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uid           []uint32               `protobuf:"varint,1,rep,packed,name=uid,proto3" json:"uid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOnlineUserIdsReply) Reset() {
	*x = GetOnlineUserIdsReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUserIdsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUserIdsReply) ProtoMessage() {}

func (x *GetOnlineUserIdsReply) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOnlineUserIdsReply.ProtoReflect.Descriptor instead.
func (*GetOnlineUserIdsReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{1}
}

func (x *GetOnlineUserIdsReply) GetUid() []uint32 {
	if x != nil {
		return x.Uid
	}
	return nil
}

type GetConnInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint32                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CustomerId    uint32                 `protobuf:"varint,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type          string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetConnInfoRequest) Reset() {
	*x = GetConnInfoRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnInfoRequest) ProtoMessage() {}

func (x *GetConnInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConnInfoRequest.ProtoReflect.Descriptor instead.
func (*GetConnInfoRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{2}
}

func (x *GetConnInfoRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetConnInfoRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *GetConnInfoRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type GetConnInfoReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Exist         bool                   `protobuf:"varint,1,opt,name=exist,proto3" json:"exist,omitempty"`
	Platform      string                 `protobuf:"bytes,2,opt,name=platform,proto3" json:"platform,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetConnInfoReply) Reset() {
	*x = GetConnInfoReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnInfoReply) ProtoMessage() {}

func (x *GetConnInfoReply) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConnInfoReply.ProtoReflect.Descriptor instead.
func (*GetConnInfoReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{3}
}

func (x *GetConnInfoReply) GetExist() bool {
	if x != nil {
		return x.Exist
	}
	return false
}

func (x *GetConnInfoReply) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

type SendUserMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MsgId         uint32                 `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendUserMessageRequest) Reset() {
	*x = SendUserMessageRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendUserMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendUserMessageRequest) ProtoMessage() {}

func (x *SendUserMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendUserMessageRequest.ProtoReflect.Descriptor instead.
func (*SendUserMessageRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{4}
}

func (x *SendUserMessageRequest) GetMsgId() uint32 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

type SendUserMessageReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendUserMessageReply) Reset() {
	*x = SendUserMessageReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendUserMessageReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendUserMessageReply) ProtoMessage() {}

func (x *SendUserMessageReply) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendUserMessageReply.ProtoReflect.Descriptor instead.
func (*SendUserMessageReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{5}
}

type SendAdminMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MsgId         uint32                 `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendAdminMessageRequest) Reset() {
	*x = SendAdminMessageRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendAdminMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendAdminMessageRequest) ProtoMessage() {}

func (x *SendAdminMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendAdminMessageRequest.ProtoReflect.Descriptor instead.
func (*SendAdminMessageRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{6}
}

func (x *SendAdminMessageRequest) GetMsgId() uint32 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

type SendAdminMessageReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendAdminMessageReply) Reset() {
	*x = SendAdminMessageReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendAdminMessageReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendAdminMessageReply) ProtoMessage() {}

func (x *SendAdminMessageReply) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendAdminMessageReply.ProtoReflect.Descriptor instead.
func (*SendAdminMessageReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{7}
}

type NoticeReadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MsgId         []uint32               `protobuf:"varint,1,rep,packed,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	UserId        uint32                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CustomerId    uint32                 `protobuf:"varint,3,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type          string                 `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeReadRequest) Reset() {
	*x = NoticeReadRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeReadRequest) ProtoMessage() {}

func (x *NoticeReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeReadRequest.ProtoReflect.Descriptor instead.
func (*NoticeReadRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{8}
}

func (x *NoticeReadRequest) GetMsgId() []uint32 {
	if x != nil {
		return x.MsgId
	}
	return nil
}

func (x *NoticeReadRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *NoticeReadRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *NoticeReadRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type NoticeReadReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeReadReply) Reset() {
	*x = NoticeReadReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeReadReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeReadReply) ProtoMessage() {}

func (x *NoticeReadReply) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeReadReply.ProtoReflect.Descriptor instead.
func (*NoticeReadReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{9}
}

var File_chat_v1_chat_proto protoreflect.FileDescriptor

var file_chat_v1_chat_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68, 0x61, 0x74, 0x22, 0x4e, 0x0a, 0x17, 0x47, 0x65,
	0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x29, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d,
	0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x62, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x44, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78,
	0x69, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x22,
	0x2f, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64,
	0x22, 0x16, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x30, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x53, 0x65,
	0x6e, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x78, 0x0a, 0x11, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x11, 0x0a,
	0x0f, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x32, 0xf2, 0x02, 0x0a, 0x04, 0x43, 0x68, 0x61, 0x74, 0x12, 0x4e, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x12, 0x1d, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3f, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4b, 0x0a, 0x0f, 0x53, 0x65,
	0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4e, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3c, 0x0a, 0x0a, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x61, 0x64, 0x12, 0x17, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74,
	0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x15, 0x5a, 0x13, 0x67, 0x66, 0x2d, 0x63, 0x68, 0x61, 0x74,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chat_v1_chat_proto_rawDescOnce sync.Once
	file_chat_v1_chat_proto_rawDescData = file_chat_v1_chat_proto_rawDesc
)

func file_chat_v1_chat_proto_rawDescGZIP() []byte {
	file_chat_v1_chat_proto_rawDescOnce.Do(func() {
		file_chat_v1_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_v1_chat_proto_rawDescData)
	})
	return file_chat_v1_chat_proto_rawDescData
}

var file_chat_v1_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_chat_v1_chat_proto_goTypes = []any{
	(*GetOnlineUserIdsRequest)(nil), // 0: chat.GetOnlineUserIdsRequest
	(*GetOnlineUserIdsReply)(nil),   // 1: chat.GetOnlineUserIdsReply
	(*GetConnInfoRequest)(nil),      // 2: chat.GetConnInfoRequest
	(*GetConnInfoReply)(nil),        // 3: chat.GetConnInfoReply
	(*SendUserMessageRequest)(nil),  // 4: chat.SendUserMessageRequest
	(*SendUserMessageReply)(nil),    // 5: chat.SendUserMessageReply
	(*SendAdminMessageRequest)(nil), // 6: chat.SendAdminMessageRequest
	(*SendAdminMessageReply)(nil),   // 7: chat.SendAdminMessageReply
	(*NoticeReadRequest)(nil),       // 8: chat.NoticeReadRequest
	(*NoticeReadReply)(nil),         // 9: chat.NoticeReadReply
}
var file_chat_v1_chat_proto_depIdxs = []int32{
	0, // 0: chat.Chat.GetOnlineUserIds:input_type -> chat.GetOnlineUserIdsRequest
	2, // 1: chat.Chat.GetConnInfo:input_type -> chat.GetConnInfoRequest
	4, // 2: chat.Chat.SendUserMessage:input_type -> chat.SendUserMessageRequest
	6, // 3: chat.Chat.SendAdminMessage:input_type -> chat.SendAdminMessageRequest
	8, // 4: chat.Chat.NoticeRead:input_type -> chat.NoticeReadRequest
	1, // 5: chat.Chat.GetOnlineUserIds:output_type -> chat.GetOnlineUserIdsReply
	3, // 6: chat.Chat.GetConnInfo:output_type -> chat.GetConnInfoReply
	5, // 7: chat.Chat.SendUserMessage:output_type -> chat.SendUserMessageReply
	7, // 8: chat.Chat.SendAdminMessage:output_type -> chat.SendAdminMessageReply
	9, // 9: chat.Chat.NoticeRead:output_type -> chat.NoticeReadReply
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_chat_v1_chat_proto_init() }
func file_chat_v1_chat_proto_init() {
	if File_chat_v1_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_v1_chat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_v1_chat_proto_goTypes,
		DependencyIndexes: file_chat_v1_chat_proto_depIdxs,
		MessageInfos:      file_chat_v1_chat_proto_msgTypes,
	}.Build()
	File_chat_v1_chat_proto = out.File
	file_chat_v1_chat_proto_rawDesc = nil
	file_chat_v1_chat_proto_goTypes = nil
	file_chat_v1_chat_proto_depIdxs = nil
}
