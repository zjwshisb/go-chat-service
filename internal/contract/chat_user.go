package contract

type IChatUser interface {
	GetUser() any
	GetPrimaryKey() uint
	GetUsername() string
	GetAvatarUrl() string
	GetCustomerId() uint
	AccessTo(user IChatUser) bool
}

type IDao interface {
}
