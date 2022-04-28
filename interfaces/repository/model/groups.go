package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Group struct {
	GroupID     uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(32)"`
	Link        string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`

	Admin []*GroupUserAdmin `gorm:"foreignKey:UserID"`
}

func (*Group) TableName() string {
	return "groups"
}

type GroupUserBelonging struct {
	UserID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	GroupID       uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	SinceYear     int       `gorm:"type:smallint(4);not null"`
	SinceSemester int       `gorm:"type:tinyint(1);not null"`
	UntilYear     int       `gorm:"type:smallint(4);not null"`
	UntilSemester int       `gorm:"type:tinyint(1);not null"`
	CreatedAt     time.Time `gorm:"precision:6"`
	UpdatedAt     time.Time `gorm:"precision:6"`

	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:UserID"`
}

func (*GroupUserBelonging) TableName() string {
	return "group_user_belongings"
}

type GroupUserAdmin struct {
	UserID  uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	GroupID uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
}

func (*GroupUserAdmin) TableName() string {
	return "group_user_admins"
}
