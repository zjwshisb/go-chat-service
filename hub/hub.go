package hub

var Hub *hub

type hub struct {
	Server *serverHub
	User *userHub
}

const (

	userLogin = "USER_LOGIN"
	userLogout = "USER_LOGOUT"
)

func Setup()  {
	server := &serverHub{
		serverClientMap: serverClientMap{
			Clients: make(map[int64]*Client),
		},
	}
	server.setup()
	user := &userHub{
		WaitingClient: &UserClientMap{
			Clients: make(map[int64]*UClient),
		},
		AcceptedClient: &UserClientMap{
			Clients: make(map[int64]*UClient),
		},
	}
	Hub = &hub{
		Server: server,
		User: user,
	}
}

