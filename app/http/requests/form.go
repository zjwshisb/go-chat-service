package requests

type AutoMessageForm struct {
	Name    string `json:"name" form:"name" binding:"required,max=32"`
	Type    string `json:"type" form:"type" binding:"required,autoMessageType"`
	Content string `json:"content" form:"content" binding:"required,max=512"`
	Title   string `json:"title" form:"title" binding:"max=32"`
	Url     string `json:"url" form:"url" binding:"max=512"`
}

type AutoRuleForm struct {
	Name      string   `json:"name" binding:"required,max=32"`
	Match     string   `json:"match" binding:"required,autoRule"`
	MatchType string   `json:"match_type" binding:"required"`
	ReplyType string   `json:"reply_type" binding:"required"`
	MessageId uint     `json:"message_id"`
	IsOpen    bool     `json:"is_open" form:"is_open"`
	Key       string   `json:"key" form:"key"`
	Sort      uint8    `json:"sort" form:"sort" binding:"required,max=128,min=0"`
	Scenes    []string `json:"scenes" form:"scenes"`
}

type AdminChatSettingForm struct {
	Background     string `json:"background" binding:"max=512"`
	IsAutoAccept   bool   `json:"is_auto_accept"`
	WelcomeContent string `json:"welcome_content" binding:"max=512"`
	OfflineContent string `json:"offline_content" binding:"max=512"`
	Name           string `json:"name" binding:"max=20"`
}
type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
