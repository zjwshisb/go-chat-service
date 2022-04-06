package contract

type FrontendUser interface {
	GetMpOpenId() string
}
type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	GetGroupId() int64
	AccessTo(user User) bool
}
