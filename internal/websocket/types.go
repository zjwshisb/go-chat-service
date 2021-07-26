package websocket

import "ws/internal/action"

type ConnManager interface {
	SendAction(act *action.Action, conn ...Conn)
	AddConn(connect Conn)
	RemoveConn(key int64)
	GetConn(key int64) (Conn, bool)
	GetAllConn() []Conn
	Login(connect Conn)
	Logout(connect Conn)
	Ping()
	Run()
}

type AnonymousConn interface {
	readMsg()
	sendMsg()
	close()
	run()
	Deliver(action *action.Action)
}
type Authentic interface {
	GetUserId() int64
}
