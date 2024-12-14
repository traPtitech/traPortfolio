package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Contest struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*Contest) TableName() string {
	return "contests"
}

type ContestLink struct {
	ContestID uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order     int       `gorm:"type:int;not null;primaryKey"`
	Link      string    `gorm:"type:text;not null"`
	// 256件とかリンク追加されても困るし小さくした方がいいか…？
	// Order uint8     `gorm:"type:tinyint(1);not null;primaryKey"`
}

func (*ContestLink) TableName() string {
	return "contest_links"
}

type ContestTeam struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	ContestID   uuid.UUID `gorm:"type:char(36);not null"`
	Name        string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:text"`
	Result      string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`

	Contest Contest `gorm:"foreignKey:ContestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*ContestTeam) TableName() string {
	return "contest_teams"
}

type ContestTeamUserBelonging struct {
	TeamID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`

	ContestTeam ContestTeam `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User        User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*ContestTeamUserBelonging) TableName() string {
	return "contest_team_user_belongings"
}

type ContestTeamLink struct {
	TeamID uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order  int       `gorm:"type:int;not null;primaryKey"`
	Link   string    `gorm:"type:text;not null"`
}

func (*ContestTeamLink) TableName() string {
	return "contest_team_links"
}
