package model

type File struct {
	Path     string `json:"path"`
	Url      string `json:"url"`
	ThumbUrl string `json:"thumb_url"`
	Storage  string `json:"storage"`
}
