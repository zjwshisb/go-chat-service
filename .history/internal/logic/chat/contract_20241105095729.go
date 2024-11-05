package chat

import (
	"gf-chat/internal/contract"
	"gf-chat/internal/model/chat"
	"github.com/gogf/gf/v2/net/ghttp"
)

type iWsConn interface {
	ReadMsg()
	SendMsg()
	Close()
	Run()
	Deliver(action *chat.Action)
	GetUserId() int
	GetUser() contract.IChatUser
	GetUuid() string
	GetPlatform() string
	GetCustomerId() int
	GetCreateTime() int64
}

// ConnContainer 管理相关方法
type connContainer interface {
	AddConn(conn iWsConn)
	GetConn(customerId int, uid int) (iWsConn, bool)
	NoticeRepeatConnect(user contract.IChatUser, newUid string)
	GetAllConn(customerId int) []iWsConn
	GetOnlineTotal(customerId int) int
	ConnExist(customerId int, uid int) bool
	Register(conn *ghttp.WebSocket, user contract.IChatUser, platform string)
	Unregister(connect iWsConn)
	RemoveConn(user contract.IChatUser)
	IsOnline(customerId int, uid int) bool
	IsLocalOnline(customerId int, uid int) bool
	GetOnlineUserIds(gid int) []int
}

type connManager interface {
	connContainer
	Run()
	Destroy()
	Ping()
	SendAction(act *chat.Action, conn ...iWsConn)
	ReceiveMessage(cm *chatConnMessage)
	GetTypes() string
	NoticeRead(customerId int, uid int, msgIds []int64)
}
