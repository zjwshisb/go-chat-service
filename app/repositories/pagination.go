package repositories

type Pagination[T any] struct {
	OriginData []T   `json:"-"`
	Total      int64 `json:"total"`
	Success    bool  `json:"success"`
	Data       interface{} `json:"data"`
}

func (p *Pagination[T]) DataFormat(f func(i T) interface{}) error {
	length := len(p.OriginData)
	data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		data[i] = f(p.OriginData[i])
	}
	p.Data = data
	return nil
}

func NewPagination[T any](data []T, total int64) *Pagination[T] {
	return &Pagination[T]{
		OriginData: data,
		Total:      total,
		Success:    true,
		Data: data,
	}
}
