package model

import (
	"github.com/gofrs/uuid"
)

type Group struct {
	GroupID uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name    string    `gorm:"type:varchar(32)"`
	Link    string    `gorm:"type:text"`
	Leader  uuid.UUID `gorm:"type:char(36);not null"`
}

type GroupUser struct {
	UserID   uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name     string    `gorm:"type:varchar(32)"`
	Realname string    `gorm:"type:varchar(32)"`
}

type GroupUserBelonging struct {
	UserID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	GroupID       uuid.UUID `gorm:"type:char(36);not null"`
	SinceYear     uint      `gorm:"type:tinyint(1);not null"`
	SinceSemester uint      `gorm:"type:tinyint(1);not null"`
	UntilYear     uint      `gorm:"type:tinyint(1);not null"`
	UntilSemester uint      `gorm:"type:tinyint(1);not null"`

	Group Group `gorm:"foreignKey:GroupID"`
}
