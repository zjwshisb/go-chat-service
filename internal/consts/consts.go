package consts

import "gf-chat/api"

const (
	MessageTypeText     = "text"
	MessageTypeImage    = "image"
	MessageTypeAudio    = "audio"
	MessageTypeVideo    = "video"
	MessageTypePdf      = "pdf"
	MessageTypeNavigate = "navigator"
	MessageTypeRate     = "rate"
	MessageTypeNotice   = "notice"
)

var MessageTypeFileTypes = []string{
	MessageTypeImage,
	MessageTypeAudio,
	MessageTypeVideo,
	MessageTypePdf,
}

var UserAllowMessageType = []api.Option{
	{
		Label: "文本",
		Value: MessageTypeText,
	},
	{
		Label: "图片",
		Value: MessageTypeImage,
	},
	{
		Label: "语音",
		Value: MessageTypeAudio,
	},
	{
		Label: "视频",
		Value: MessageTypeVideo,
	},
	{
		Label: "PDF",
		Value: MessageTypePdf,
	},
	{
		Label: "导航卡片",
		Value: MessageTypeNavigate,
	},
}

const (
	ChatSessionTypeNormal   = 0
	ChatSessionTypeTransfer = 1
)

const WebsocketPlatformWeb = "web"
const WebsocketPlatformH5 = "h5"
const WebsocketPlatformMp = "weapp"

const (
	ActionReceipt          = "receipt"
	ActionPing             = "ping"
	ActionUserOnLine       = "user-online"
	ActionUserOffLine      = "user-offline"
	ActionWaitingUser      = "waiting-users"
	ActionWaitingUserCount = "waiting-user-count"
	ActionAdmins           = "admins"
	ActionSendMessage      = "send-message"
	ActionReceiveMessage   = "receive-message"
	ActionOtherLogin       = "other-login"
	ActionMoreThanOne      = "more-than-one"
	ActionUserTransfer     = "user-transfer"
	ActionErrorMessage     = "error-message"
	ActionRead             = "read"
	ActionUserRate         = "user-rate"
)

const (
	AutoRuleMatchTypeAll  = "all"
	AutoRuleMatchTypePart = "part"

	AutoRuleMatchEnter           = "enter"
	AutoRuleMatchAdminAllOffLine = "u-offline"

	AutoRuleReplyTypeMessage  = "message"
	AutoRuleReplyTypeTransfer = "transfer"

	AutoRuleSceneNotAccepted  = "not-accepted"
	AutoRuleSceneAdminOnline  = "admin-online"
	AutoRuleSceneAdminOffline = "admin-offline"
)

const (
	ChatSessionStatusWait   = "wait"
	ChatSessionStatusCancel = "cancel"
	ChatSessionStatusAccept = "accept"
	ChatSessionStatusClose  = "close"
)

const (
	MessageSourceUser   = 0
	MessageSourceAdmin  = 1
	MessageSourceSystem = 2
)

const (
	ChatSettingTypeImage  = "image"
	ChatSettingTypeText   = "text"
	ChatSettingTypeSelect = "select"

	ChatSettingIsAutoTransfer = "is-auto-transfer"
	ChatSettingMinuteToBreak  = "minute-to-break"
	ChatSettingSystemName     = "system-name"
	ChatSettingSystemAvatar   = "system-avatar"
	ChatSettingShowQueue      = "show-queue"
	ChatSettingShowRead       = "show-read"
)

const (
	StorageQiniu = "qiniu"
	StorageLocal = "local"
)

const (
	FileTypeDir   = "dir"
	FileTypeImage = "image"
	FileTypeVideo = "video"
	FileTypeAudio = "audio"
)
