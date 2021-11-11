package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*Project) TableName() string {
	return "projects"
}

type ProjectMember struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	ProjectID uuid.UUID `gorm:"type:char(36);not null"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	Since     time.Time `gorm:"precision:6"`
	Until     time.Time `gorm:"precision:6"`

	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
}

func (*ProjectMember) TableName() string {
	return "project_members"
}
