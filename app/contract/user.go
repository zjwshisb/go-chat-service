package contract

type FrontendUser interface {
	GetMpOpenId() string
}
type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	GetGroupId() int64
	GetReqId() string
	AccessTo(user User) bool
}
