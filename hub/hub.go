package hub

var Hub *hub

type hub struct {
	Server *serverHub
	User *userHub
}

const (
	serverLogin = "LOGIN"
	serverLogout = "LOGOUT"
	userLogin = "USER_LOGIN"
	userLogout = "USER_LOGOUT"
)

func Setup()  {
	server := &serverHub{
		Clients: make(map[int64]*Client),
	}
	server.setup()
	user := &userHub{
		Clients: make(map[int64]*UClient),
		Waiting: make(map[int64]*UClient),
	}
	Hub = &hub{
		Server: server,
		User: user,
	}
}

