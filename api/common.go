package api

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
