package request

type ConnectionExistRequest struct {
	Uid int64
}

type ConnectionGroupRequest struct {
	GroupId int64
	Types   string
}

type ConnectionAllCountRequest struct {
	Types string
}
