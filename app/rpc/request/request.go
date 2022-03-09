package request

type ExistRequest struct {
	Uid int64
}

type GroupRequest struct {
	GroupId int64
	Types string
}

