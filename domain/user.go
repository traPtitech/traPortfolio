package domain

type User struct {
	ID          uint
	Name        string
	DisplayName string
	Status      UserAccountStatus
}

type UserAccountStatus int
