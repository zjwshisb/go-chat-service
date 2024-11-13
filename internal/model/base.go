package model

type Paginator[I any] struct {
	Items    []I
	Total    int
	IsSimple bool
}

type ImageFiled struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

type Option struct {
	Value any    `json:"value"`
	Label string `json:"label"`
}

type QueryInput struct {
	Size int
	Page int
}

func (p *QueryInput) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *QueryInput) GetSize() int {
	if p.Size <= 0 || p.Size >= 200 {
		p.Size = 20
	}
	return p.Size
}
