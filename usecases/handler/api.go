package handler

type API struct {
	Ping *PingHandler
	User *UserHandler
}

func NewAPI(ping *PingHandler, user *UserHandler) API {
	return API{
		Ping: ping,
		User: user,
	}
}
