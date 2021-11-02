package handler

type API struct {
	Ping    *PingHandler
	User    *UserHandler
	Project *ProjectHandler
	Event   *EventHandler
	Contest *ContestHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, project *ProjectHandler, event *EventHandler, contest *ContestHandler) API {
	return API{
		Ping:    ping,
		User:    user,
		Project: project,
		Event:   event,
		Contest: contest,
	}
}
