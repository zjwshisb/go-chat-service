package json

type Options struct {
	Value interface{} `json:"value"`
	Label string   `json:"label"`
}
type Setting struct {
	Name string `json:"name"`
	Title string `json:"title"`
	Value string `json:"value"`
	Options map[string]string `json:"options"`
}