package chat

import (
	api "gf-chat/api/v1/backend"
	"github.com/gorilla/websocket"
)

type iWsConn interface {
	ReadMsg()
	SendMsg()
	Close()
	Run()
	Deliver(action *api.ChatAction)
	GetUserId() uint
	GetUser() IChatUser
	GetUuid() string
	GetPlatform() string
	GetCustomerId() uint
	GetCreateTime() int64
}

// ConnContainer 管理相关方法
type connContainer interface {
	AddConn(conn iWsConn)
	GetConn(customerId uint, uid uint) (iWsConn, bool)
	NoticeRepeatConnect(user IChatUser, newUid string)
	GetAllConn(customerId uint) []iWsConn
	GetOnlineTotal(customerId uint) uint
	ConnExist(customerId uint, uid uint) bool
	Register(conn *websocket.Conn, user IChatUser, platform string)
	Unregister(connect iWsConn)
	RemoveConn(user IChatUser)
	IsOnline(customerId uint, uid uint) bool
	IsLocalOnline(customerId uint, uid uint) bool
	GetOnlineUserIds(gid uint) []uint
}

type connManager interface {
	connContainer
	Run()
	Destroy()
	Ping()
	SendAction(act *api.ChatAction, conn ...iWsConn)
	ReceiveMessage(cm *chatConnMessage)
	GetTypes() string
	NoticeRead(customerId uint, uid uint, msgIds []uint)
}
