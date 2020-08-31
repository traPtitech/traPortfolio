package domain

type User struct {
	ID          uint
	name        string
	displayName string
	state       UserAccountStatus
}

type UserAccountStatus int
