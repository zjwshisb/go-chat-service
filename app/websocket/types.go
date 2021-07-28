package websocket


type ConnManager interface {
	SendAction(act *Action, conn ...Conn)
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
	Deliver(action *Action)
}

type Authentic interface {
	GetUserId() int64
}

type Handle func(i ...interface{})

