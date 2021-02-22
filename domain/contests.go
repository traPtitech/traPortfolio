package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Contest struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(32);not null;unique"`
	Description string    `gorm:"type:text;not null"`
	Link        string    `gorm:"type:text;not null"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

type ContestTeam struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	ContestID   uuid.UUID `gorm:"type:char(36);not null;unique"`
	Name        string    `gorm:"type:varchar(32);not null;unique"`
	Description string    `gorm:"type:text;not null"`
	Result      string    `gorm:"type:text;not null"`
	Link        string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}
