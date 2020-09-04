package input

import "github.com/traPtitech/traPortfolio/domain"

type GetUser struct {
	Id int
}

type AddUser struct {
	User domain.User
}

type UpdateUser struct {
	User domain.User
}

type DeleteUser struct {
	Id int
}
