package domain

import "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	RealName string
}
