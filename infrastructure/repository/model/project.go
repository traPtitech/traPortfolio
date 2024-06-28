package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID            uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name          string    `gorm:"type:varchar(32)"`
	Description   string    `gorm:"type:text"`
	Link          string    `gorm:"type:text"`
	SinceYear     int       `gorm:"type:smallint(4);not null"`
	SinceSemester int       `gorm:"type:tinyint(1);not null"`
	UntilYear     int       `gorm:"type:smallint(4);not null"`
	UntilSemester int       `gorm:"type:tinyint(1);not null"`
	CreatedAt     time.Time `gorm:"precision:6"`
	UpdatedAt     time.Time `gorm:"precision:6"`
}

func (*Project) TableName() string {
	return "projects"
}

type ProjectMember struct {
	ProjectID     uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	UserID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	SinceYear     int       `gorm:"type:smallint(4);not null"`
	SinceSemester int       `gorm:"type:tinyint(1);not null"`
	UntilYear     int       `gorm:"type:smallint(4);not null"`
	UntilSemester int       `gorm:"type:tinyint(1);not null"`
	CreatedAt     time.Time `gorm:"precision:6"`
	UpdatedAt     time.Time `gorm:"precision:6"`

	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*ProjectMember) TableName() string {
	return "project_members"
}
