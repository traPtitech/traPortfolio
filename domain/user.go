package domain

import "github.com/gofrs/uuid"

type User struct {
	Id       uuid.UUID
	Name     string
	RealName string
}
