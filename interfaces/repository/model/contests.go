package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Contest struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

type ContestTeam struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text"`
	Result      string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

type ContestContestTeamBelonging struct {
	ContestID uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	TeamID    uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

type ContestTeamUserBelonging struct {
	TeamID    uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}
