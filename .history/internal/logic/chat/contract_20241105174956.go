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
	GetUserId() uint
	GetUser() contract.IChatUser
	GetUuid() string
	GetPlatform() string
	GetCustomerId() uint
	GetCreateTime() int64
}

// ConnContainer 管理相关方法
type connContainer interface {
	AddConn(conn iWsConn)
	GetConn(customerId uint, uid uint) (iWsConn, bool)
	NoticeRepeatConnect(user contract.IChatUser, newUid string)
	GetAllConn(customerId uint) []iWsConn
	GetOnlineTotal(customerId uint) uint
	ConnExist(customerId uint, uid uint) bool
	Register(conn *ghttp.WebSocket, user contract.IChatUser, platform string)
	Unregister(connect iWsConn)
	RemoveConn(user contract.IChatUser)
	IsOnline(customerId uint, uid uint) bool
	IsLocalOnline(customerId uint, uid uint) bool
	GetOnlineUserIds(gid uint) []uint
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
