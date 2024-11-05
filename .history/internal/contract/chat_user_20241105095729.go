package contract

type IChatUser interface {
	GetUser() any
	GetPrimaryKey() int
	GetUsername() string
	GetAvatarUrl() string
	GetCustomerId() int
	AccessTo(user IChatUser) bool
}
