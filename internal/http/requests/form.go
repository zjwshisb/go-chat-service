package requests

type AutoMessageForm struct {
	Name string `json:"name" form:"name" binding:"required,max=32"`
	Type string `json:"type" form:"type" binding:"required,autoMessageType"`
	Content string `json:"content" form:"content" binding:"required,max=512"`
	Title string `json:"title" form:"title" binding:"max=32"`
	Url string `json:"url" form:"url" binding:"max=512"`
}

type AutoRuleForm struct {
	Name string `json:"name" binding:"required,max=32"`
	Match string `json:"match" binding:"required,autoRule"`
	MatchType string `json:"match_type" binding:"required"`
	ReplyType string `json:"reply_type" binding:"required"`
	MessageId uint `json:"message_id"`
	IsOpen bool `json:"is_open" form:"is_open"`
	Sort uint8 `json:"sort" form:"sort" binding:"required,max=128,min=0"`
}