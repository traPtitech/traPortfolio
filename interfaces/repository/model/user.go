package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Description string    `gorm:"type:text;not null"`
	Check       bool      `gorm:"type:boolean;not null;default:false"`
	Name        string    `gorm:"type:varchar(32);not null;unique"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*User) TableName() string {
	return "users"
}

type Account struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Type      uint      `gorm:"type:tinyint(1);not null"`
	Name      string    `gorm:"type:varchar(32)"`
	URL       string    `gorm:"type:text"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;unique"`
	Check     bool      `gorm:"type:boolean;not null;default:false"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}
