package hub

var Hub *hub

type hub struct {
	Server *serverHub
	User *userHub
}

func Setup()  {
	server := &serverHub{
		Clients: make(map[int64]*Client),
	}
	user := &userHub{
		Clients: make(map[int64]*UClient),
		Waiting: make([]*UClient, 0),
	}
	Hub = &hub{
		Server: server,
		User: user,
	}
}

