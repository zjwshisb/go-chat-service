package response

type NilResponse struct {
}

type OnlineResponse struct {
	Data bool
}
type IdsResponse struct {
	Data []int64
}

type CountResponse struct {
	Data int64
}
type SendMessageResponse struct {
}
