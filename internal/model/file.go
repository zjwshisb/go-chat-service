package model

type File struct {
	Id       uint   `json:"id"`
	Path     string `json:"path"`
	Url      string `json:"url"`
	ThumbUrl string `json:"thumb_url"`
	Type     string `json:"type"`
}

type SaveFileOutput struct {
	Path string
	Disk string
}
