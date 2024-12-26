package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
)

type User struct {
	ID          uuid.UUID        `gorm:"type:char(36);not null;primaryKey"`
	Description string           `gorm:"type:text;not null"`
	Check       bool             `gorm:"type:boolean;not null;default:false"`
	Name        string           `gorm:"type:varchar(32);not null;unique"`
	DisplayName string           `gorm:"type:varchar(32)"`
	State       domain.TraQState `gorm:"type:tinyint(1);not null"`
	CreatedAt   time.Time        `gorm:"precision:6"`
	UpdatedAt   time.Time        `gorm:"precision:6"`

	Accounts []*Account `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*User) TableName() string {
	return "users"
}

type Account struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Type      uint8     `gorm:"type:tinyint(1);not null"`
	Name      string    `gorm:"type:varchar(256)"`
	URL       string    `gorm:"type:text"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	Check     bool      `gorm:"type:boolean;not null;default:false"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (*Account) TableName() string {
	return "accounts"
}
