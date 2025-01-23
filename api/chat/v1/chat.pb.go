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

type NilReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NilReply) Reset() {
	*x = NilReply{}
	mi := &file_chat_v1_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NilReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NilReply) ProtoMessage() {}

func (x *NilReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use NilReply.ProtoReflect.Descriptor instead.
func (*NilReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{0}
}

type GetOnlineUserIdsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerId    uint32                 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOnlineUserIdsRequest) Reset() {
	*x = GetOnlineUserIdsRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUserIdsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUserIdsRequest) ProtoMessage() {}

func (x *GetOnlineUserIdsRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetOnlineUserIdsRequest.ProtoReflect.Descriptor instead.
func (*GetOnlineUserIdsRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{1}
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
	mi := &file_chat_v1_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUserIdsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUserIdsReply) ProtoMessage() {}

func (x *GetOnlineUserIdsReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetOnlineUserIdsReply.ProtoReflect.Descriptor instead.
func (*GetOnlineUserIdsReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{2}
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
	mi := &file_chat_v1_chat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnInfoRequest) ProtoMessage() {}

func (x *GetConnInfoRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetConnInfoRequest.ProtoReflect.Descriptor instead.
func (*GetConnInfoRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{3}
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
	mi := &file_chat_v1_chat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetConnInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConnInfoReply) ProtoMessage() {}

func (x *GetConnInfoReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetConnInfoReply.ProtoReflect.Descriptor instead.
func (*GetConnInfoReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{4}
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

type SendMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MsgId         uint32                 `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{5}
}

func (x *SendMessageRequest) GetMsgId() uint32 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *SendMessageRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
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
	mi := &file_chat_v1_chat_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeReadRequest) ProtoMessage() {}

func (x *NoticeReadRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use NoticeReadRequest.ProtoReflect.Descriptor instead.
func (*NoticeReadRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{6}
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
	mi := &file_chat_v1_chat_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeReadReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeReadReply) ProtoMessage() {}

func (x *NoticeReadReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use NoticeReadReply.ProtoReflect.Descriptor instead.
func (*NoticeReadReply) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{7}
}

type UpdateAdminSettingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateAdminSettingRequest) Reset() {
	*x = UpdateAdminSettingRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateAdminSettingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAdminSettingRequest) ProtoMessage() {}

func (x *UpdateAdminSettingRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use UpdateAdminSettingRequest.ProtoReflect.Descriptor instead.
func (*UpdateAdminSettingRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateAdminSettingRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type BroadcastWaitingUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerId    uint32                 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BroadcastWaitingUserRequest) Reset() {
	*x = BroadcastWaitingUserRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BroadcastWaitingUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastWaitingUserRequest) ProtoMessage() {}

func (x *BroadcastWaitingUserRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use BroadcastWaitingUserRequest.ProtoReflect.Descriptor instead.
func (*BroadcastWaitingUserRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{9}
}

func (x *BroadcastWaitingUserRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type NoticeTransferRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AdminId       uint32                 `protobuf:"varint,1,opt,name=admin_id,json=adminId,proto3" json:"admin_id,omitempty"`
	CustomerId    uint32                 `protobuf:"varint,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeTransferRequest) Reset() {
	*x = NoticeTransferRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeTransferRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeTransferRequest) ProtoMessage() {}

func (x *NoticeTransferRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeTransferRequest.ProtoReflect.Descriptor instead.
func (*NoticeTransferRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{10}
}

func (x *NoticeTransferRequest) GetAdminId() uint32 {
	if x != nil {
		return x.AdminId
	}
	return 0
}

func (x *NoticeTransferRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type NoticeUserOnlineRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint32                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Platform      string                 `protobuf:"bytes,2,opt,name=platform,proto3" json:"platform,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeUserOnlineRequest) Reset() {
	*x = NoticeUserOnlineRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeUserOnlineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeUserOnlineRequest) ProtoMessage() {}

func (x *NoticeUserOnlineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeUserOnlineRequest.ProtoReflect.Descriptor instead.
func (*NoticeUserOnlineRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{11}
}

func (x *NoticeUserOnlineRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *NoticeUserOnlineRequest) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

type NoticeUserOfflineRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint32                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeUserOfflineRequest) Reset() {
	*x = NoticeUserOfflineRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeUserOfflineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeUserOfflineRequest) ProtoMessage() {}

func (x *NoticeUserOfflineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeUserOfflineRequest.ProtoReflect.Descriptor instead.
func (*NoticeUserOfflineRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{12}
}

func (x *NoticeUserOfflineRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type BroadcastOnlineAdminsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerId    uint32                 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BroadcastOnlineAdminsRequest) Reset() {
	*x = BroadcastOnlineAdminsRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BroadcastOnlineAdminsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastOnlineAdminsRequest) ProtoMessage() {}

func (x *BroadcastOnlineAdminsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastOnlineAdminsRequest.ProtoReflect.Descriptor instead.
func (*BroadcastOnlineAdminsRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{13}
}

func (x *BroadcastOnlineAdminsRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type BroadcastQueueLocationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CustomerId    uint32                 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BroadcastQueueLocationRequest) Reset() {
	*x = BroadcastQueueLocationRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BroadcastQueueLocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastQueueLocationRequest) ProtoMessage() {}

func (x *BroadcastQueueLocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastQueueLocationRequest.ProtoReflect.Descriptor instead.
func (*BroadcastQueueLocationRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{14}
}

func (x *BroadcastQueueLocationRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type NoticeRepeatConnectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint32                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CustomerId    uint32                 `protobuf:"varint,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type          string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	NewUid        string                 `protobuf:"bytes,4,opt,name=new_uid,json=newUid,proto3" json:"new_uid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeRepeatConnectRequest) Reset() {
	*x = NoticeRepeatConnectRequest{}
	mi := &file_chat_v1_chat_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeRepeatConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeRepeatConnectRequest) ProtoMessage() {}

func (x *NoticeRepeatConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_v1_chat_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeRepeatConnectRequest.ProtoReflect.Descriptor instead.
func (*NoticeRepeatConnectRequest) Descriptor() ([]byte, []int) {
	return file_chat_v1_chat_proto_rawDescGZIP(), []int{15}
}

func (x *NoticeRepeatConnectRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *NoticeRepeatConnectRequest) GetCustomerId() uint32 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *NoticeRepeatConnectRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *NoticeRepeatConnectRequest) GetNewUid() string {
	if x != nil {
		return x.NewUid
	}
	return ""
}

var File_chat_v1_chat_proto protoreflect.FileDescriptor

var file_chat_v1_chat_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68, 0x61, 0x74, 0x22, 0x0a, 0x0a, 0x08, 0x4e, 0x69,
	0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x4e, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x29, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x22, 0x62, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x44, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x69,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x22, 0x3f, 0x0a, 0x12, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x78, 0x0a, 0x11,
	0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x2b, 0x0a, 0x19, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x1b, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x57, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x22, 0x53, 0x0a, 0x15, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4e, 0x0a, 0x17, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x22, 0x33, 0x0a, 0x18, 0x4e,
	0x6f, 0x74, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x3f, 0x0a, 0x1c, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x40, 0x0a, 0x1d, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x83, 0x01, 0x0a, 0x1a, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x6e, 0x65, 0x77, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x65, 0x77, 0x55, 0x69, 0x64, 0x32, 0xcc, 0x06, 0x0a, 0x04, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x4e, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x12, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65,
	0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x3f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x37, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x18, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3c, 0x0a, 0x0a,
	0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64, 0x12, 0x17, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x12, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x1f, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x49, 0x0a, 0x14, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x57, 0x61,
	0x69, 0x74, 0x69, 0x6e, 0x67, 0x55, 0x73, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x57, 0x61, 0x69, 0x74, 0x69, 0x6e,
	0x67, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4b, 0x0a, 0x15,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x73, 0x12, 0x22, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x42, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4d, 0x0a, 0x16, 0x42, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x23, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3d, 0x0a, 0x0e, 0x4e, 0x6f, 0x74, 0x69,
	0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e,
	0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x41, 0x0a, 0x10, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x1d, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x43, 0x0a, 0x11, 0x4e, 0x6f,
	0x74, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x12,
	0x1e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x47, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x20, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x4e, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x15, 0x5a, 0x13, 0x67, 0x66, 0x2d, 0x63,
	0x68, 0x61, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_chat_v1_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_chat_v1_chat_proto_goTypes = []any{
	(*NilReply)(nil),                      // 0: chat.NilReply
	(*GetOnlineUserIdsRequest)(nil),       // 1: chat.GetOnlineUserIdsRequest
	(*GetOnlineUserIdsReply)(nil),         // 2: chat.GetOnlineUserIdsReply
	(*GetConnInfoRequest)(nil),            // 3: chat.GetConnInfoRequest
	(*GetConnInfoReply)(nil),              // 4: chat.GetConnInfoReply
	(*SendMessageRequest)(nil),            // 5: chat.SendMessageRequest
	(*NoticeReadRequest)(nil),             // 6: chat.NoticeReadRequest
	(*NoticeReadReply)(nil),               // 7: chat.NoticeReadReply
	(*UpdateAdminSettingRequest)(nil),     // 8: chat.UpdateAdminSettingRequest
	(*BroadcastWaitingUserRequest)(nil),   // 9: chat.BroadcastWaitingUserRequest
	(*NoticeTransferRequest)(nil),         // 10: chat.NoticeTransferRequest
	(*NoticeUserOnlineRequest)(nil),       // 11: chat.NoticeUserOnlineRequest
	(*NoticeUserOfflineRequest)(nil),      // 12: chat.NoticeUserOfflineRequest
	(*BroadcastOnlineAdminsRequest)(nil),  // 13: chat.BroadcastOnlineAdminsRequest
	(*BroadcastQueueLocationRequest)(nil), // 14: chat.BroadcastQueueLocationRequest
	(*NoticeRepeatConnectRequest)(nil),    // 15: chat.NoticeRepeatConnectRequest
}
var file_chat_v1_chat_proto_depIdxs = []int32{
	1,  // 0: chat.Chat.GetOnlineUserIds:input_type -> chat.GetOnlineUserIdsRequest
	3,  // 1: chat.Chat.GetConnInfo:input_type -> chat.GetConnInfoRequest
	5,  // 2: chat.Chat.SendMessage:input_type -> chat.SendMessageRequest
	6,  // 3: chat.Chat.NoticeRead:input_type -> chat.NoticeReadRequest
	8,  // 4: chat.Chat.UpdateAdminSetting:input_type -> chat.UpdateAdminSettingRequest
	9,  // 5: chat.Chat.BroadcastWaitingUser:input_type -> chat.BroadcastWaitingUserRequest
	13, // 6: chat.Chat.BroadcastOnlineAdmins:input_type -> chat.BroadcastOnlineAdminsRequest
	14, // 7: chat.Chat.BroadcastQueueLocation:input_type -> chat.BroadcastQueueLocationRequest
	10, // 8: chat.Chat.NoticeTransfer:input_type -> chat.NoticeTransferRequest
	11, // 9: chat.Chat.NoticeUserOnline:input_type -> chat.NoticeUserOnlineRequest
	12, // 10: chat.Chat.NoticeUserOffline:input_type -> chat.NoticeUserOfflineRequest
	15, // 11: chat.Chat.NoticeRepeatConnect:input_type -> chat.NoticeRepeatConnectRequest
	2,  // 12: chat.Chat.GetOnlineUserIds:output_type -> chat.GetOnlineUserIdsReply
	4,  // 13: chat.Chat.GetConnInfo:output_type -> chat.GetConnInfoReply
	0,  // 14: chat.Chat.SendMessage:output_type -> chat.NilReply
	7,  // 15: chat.Chat.NoticeRead:output_type -> chat.NoticeReadReply
	0,  // 16: chat.Chat.UpdateAdminSetting:output_type -> chat.NilReply
	0,  // 17: chat.Chat.BroadcastWaitingUser:output_type -> chat.NilReply
	0,  // 18: chat.Chat.BroadcastOnlineAdmins:output_type -> chat.NilReply
	0,  // 19: chat.Chat.BroadcastQueueLocation:output_type -> chat.NilReply
	0,  // 20: chat.Chat.NoticeTransfer:output_type -> chat.NilReply
	0,  // 21: chat.Chat.NoticeUserOnline:output_type -> chat.NilReply
	0,  // 22: chat.Chat.NoticeUserOffline:output_type -> chat.NilReply
	0,  // 23: chat.Chat.NoticeRepeatConnect:output_type -> chat.NilReply
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
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
			NumMessages:   16,
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