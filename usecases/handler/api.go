package handler

type API struct {
	Ping  *PingHandler
	User  *UserHandler
	Event *EventHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, event *EventHandler) API {
	return API{
		Ping:  ping,
		User:  user,
		Event: event,
	}
}
