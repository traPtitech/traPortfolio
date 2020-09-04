package controllers

type API struct {
	Ping *PingController
	User *UserController
}

// func NewAPI(ping *PingController, user *UserController) API {
// 	return API{
// 		Ping: ping,
// 		User: user,
// 	}
// }
func NewAPI(ping *PingController) API {
	return API{
		Ping: ping,
	}
}
