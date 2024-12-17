package v1

import "github.com/gogf/gf/v2/os/gtime"

type ListRes[T any] struct {
	NormalRes[[]T]
	Total int `json:"total"`
}

type NilRes = NormalRes[any]

type Paginate struct {
	PageSize int `d:"20" json:"pageSize" v:"max:100"`
	Current  int `d:"1" dc:"页码" json:"current"`
}

type OptionRes struct {
	NormalRes[[]Option]
}

type Option struct {
	Value any    `json:"value"`
	Label string `json:"label"`
}

type NormalRes[T any] struct {
	Code    int  `json:"code"    dc:"Error code"`
	Data    T    `json:"data"    dc:"Result data for certain request according API definition"`
	Success bool `json:"success" dc:"Is Success"`
}

type FailRes struct {
	Code    int    `json:"code"    dc:"Error code"`
	Success bool   `json:"success" dc:"Is Success"`
	Message string `json:"message" dc:"错误消息"`
}

func NewOptionResp(options []Option) *OptionRes {
	return &OptionRes{
		NormalRes: NormalRes[[]Option]{
			Code:    0,
			Success: true,
			Data:    options,
		},
	}
}

func NewFailResp(message string, code int) *FailRes {
	return &FailRes{
		Code:    code,
		Success: false,
		Message: message,
	}
}

func NewListResp[T any](items []T, total int) *ListRes[T] {
	return &ListRes[T]{
		NormalRes: NormalRes[[]T]{
			Success: true,
			Data:    items,
			Code:    0,
		},
		Total: total,
	}
}
func NewNilResp() *NilRes {
	return &NilRes{
		Success: true,
		Data:    nil,
		Code:    0,
	}
}
func NewResp[T any](data T) *NormalRes[T] {
	return &NormalRes[T]{
		Success: true,
		Data:    data,
		Code:    0,
	}
}

type ChatMessage struct {
	Id         uint        `json:"id"`
	UserId     uint        `json:"user_id"`
	AdminId    uint        `json:"admin_id"`
	AdminName  string      `json:"admin_name"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	ReceivedAT *gtime.Time `json:"received_at"`
	Source     uint        `json:"source"`
	ReqId      string      `json:"req_id"`
	IsSuccess  bool        `json:"is_success"`
	IsRead     bool        `json:"is_read"`
	Avatar     string      `json:"avatar"`
	Username   string      `json:"username"`
}

type ChatAction struct {
	Data   any    `json:"data"`
	Time   int64  `json:"time"`
	Action string `json:"action"`
}
