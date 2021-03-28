package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text;not null"`
	Start       time.Time `gorm:"datetime(6)"`
	End         time.Time `gorm:"datetime(6)"`
}
