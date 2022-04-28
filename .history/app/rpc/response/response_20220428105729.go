package response

type ConnectionExistResponse struct {
	Exists bool
}

type ConnectionTotalResponse struct {
	Total int64
}

type ConnectionIdsResponse struct {
	Ids []int64
}

type ConnectionAllCountResponse struct {
	Total int64
}
