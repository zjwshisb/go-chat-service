package request

type NormalRequest struct {
	Types   string
	GroupId int64
}

type OnlineRequest struct {
	Id    int64
	Types string
}

type IdsRequest struct {
	Types   string
	GroupId int64
}

type SendMessageRequest struct {
	Id int64
}

type RepeatConnectRequest struct {
	Types string
	Id    int64
}
