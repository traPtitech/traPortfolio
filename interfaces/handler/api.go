package handler

type API struct {
	Ping    *PingHandler
	User    *UserHandler
	Event   *EventHandler
	Contest *ContestHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, event *EventHandler, contest *ContestHandler) API {
	return API{
		Ping:    ping,
		User:    user,
		Event:   event,
		Contest: contest,
	}
}
