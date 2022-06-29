package request

type NormalRequest struct {
	Types   string
	GroupId int64
}

type TypeRequest struct {
	Types string
}

type GroupRequest struct {
	GroupId int64
}

type OnlineRequest struct {
	Id    int64
	Types string
}

type IdRequest struct {
	Id int64
}

type IdsRequest struct {
	Types   string
	GroupId int64
}

type SendMessageRequest struct {
	Id int64
}

type RepeatConnectRequest struct {
	Types   string
	Id      int64
	NewUuid string
}
