package resource

import "time"

type Options struct {
	Value interface{} `json:"value"`
	Label string   `json:"label"`
}

type Line struct {
	Category string `json:"category"`
	Value int `json:"value"`
	Label interface{} `json:"label"`
}

type Admin struct {
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	Online        bool   `json:"online"`
	Id            int64  `json:"id"`
	AcceptedCount int  `json:"accepted_count"`
}

type AutoMessage struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RulesCount int       `json:"rules_count"`
}

type AutoRule struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Match       string       `json:"match"`
	MatchType   string       `json:"match_type"`
	ReplyType   string       `json:"reply_type"`
	MessageId   uint         `json:"message_id"`
	Key         string       `gorm:"key" json:"key"`
	Sort        uint8        `json:"sort"`
	IsOpen      bool         `json:"is_open"`
	Count       uint         `json:"count"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	EventLabel  string       `json:"event_label"`
	Message     *AutoMessage `json:"message"`
	Scenes      []string     `json:"scenes"`
	ScenesLabel string       `json:"scenes_label"`
}

type ChatSession struct {
	Id         uint64 `json:"id"`
	UserId     int64  `json:"-"`
	QueriedAt  int64  `json:"queried_at"`
	AcceptedAt int64  `json:"accepted_at"`
	BrokeAt    int64  `json:"broke_at"`
	CanceledAt int64 `json:"canceled_at"`
	AdminId    int64  `json:"admin_id"`
	UserName   string `json:"user_name"`
	AdminName  string `json:"admin_name"`
	TypeLabel  string `json:"type_label"`
	Status string `json:"status"`
}


type SimpleMessage struct {
	Type string `json:"type"`
	Time int64 `json:"time"`
	Content string `json:"content"`
}

type Message struct {
	Id         uint64 `json:"id"`
	UserId     int64  `json:"user_id"`
	AdminId    int64  `json:"admin_id"`
	AdminName  string `json:"admin_name"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	ReceivedAT int64  `json:"received_at"`
	Source     int8   `json:"source"`
	ReqId      string  `json:"req_id"`
	IsSuccess  bool   `json:"is_success"`
	IsRead     bool   `json:"is_read"`
	Avatar     string `json:"avatar"`
}


type WaitingChatSession struct {
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	UserId           int64  `json:"id"`
	LastTime int64 `json:"last_time"`
	Messages []*SimpleMessage `json:"messages"`
	MessageCount int    `json:"message_count"`
	Description  string `json:"description"`
	SessionId   uint64 `json:"session_id"`
}

type ChatTransfer struct {
	Id            int64      `json:"id"`
	SessionId     uint64     `json:"session_id"`
	UserId        int64      `json:"user_id"`
	Remark        string     `json:"remark"`
	FromAdminName string     `json:"from_admin_name"`
	ToAdminName   string     `json:"to_admin_name"`
	Username      string     `json:"username"`
	CreatedAt     *time.Time `json:"created_at"`
	AcceptedAt    *time.Time `json:"accepted_at"`
	CanceledAt    *time.Time `json:"canceled_at"`
}

type ChatSetting struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Title string `json:"title"`
	Value string `json:"value"`
	Options []map[string]string `json:"options"`
}


type User struct {
	ID           int64          `json:"id"`
	Username     string         `json:"username"`
	LastChatTime int64          `json:"last_chat_time"`
	Disabled     bool           `json:"disabled"`
	Online       bool           `json:"online"`
	Messages     []*Message `json:"messages"`
	Unread       int            `json:"unread"`
	Avatar string `json:"avatar"`
}
