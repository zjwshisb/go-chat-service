package json

type Options struct {
	Value interface{} `json:"value"`
	Label string   `json:"label"`
}

type Line struct {
	Category string `json:"category"`
	Value int `json:"value"`
	Label interface{} `json:"label"`
}
