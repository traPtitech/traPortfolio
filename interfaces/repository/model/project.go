package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

type ProjectMember struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	ProjectID uuid.UUID `gorm:"type:char(36);not null"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	Since     time.Time `gorm:"precision:6"`
	Until     time.Time `gorm:"precision:6"`
}

type ProjectDetail struct {
	ID          uuid.UUID
	Name        string
	Link        string
	Description string
	Members     []*ProjectMemberDetail
	Since       time.Time
	Until       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectMemberDetail struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Name      string
	RealName  string
	Since     time.Time
	Until     time.Time
}
