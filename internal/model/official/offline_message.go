package official

type Message struct {
	Miniprogram *Miniprogram `json:"miniprogram"`
	Data        any          `json:"data"`
	Url         string       `json:"url"`
	TemplateId  string       `json:"template_id"`
	ToUser      string       `json:"touser"`
}

type Miniprogram struct {
	Appid    string
	Pagepath string
}

type Offline struct {
	First    string `json:"first"`
	Remark   string `json:"remark"`
	Keyword1 string `json:"keyword1"` //客户名称
	Keyword2 string `json:"keyword2"` //客户标签
	Keyword3 string `json:"keyword3"` //消息来源
	Keyword4 string `json:"keyword4"` //咨询时间
}
