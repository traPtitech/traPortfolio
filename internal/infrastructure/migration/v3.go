// Package migration migrate current struct
package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// v3 ユーザーアカウントのprPermitted属性廃止
func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(db *gorm.DB) error {
			if err := db.Migrator().DropColumn(v3Account{}, "check"); err != nil {
				return err
			}

			return db.
				Table("portfolio").
				Error
		},
	}
}

type v3Account struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Type      uint8     `gorm:"type:tinyint(1);not null"`
	Name      string    `gorm:"type:varchar(256)"`
	URL       string    `gorm:"type:text"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	Check     bool      `gorm:"type:boolean;not null;default:false"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (*v3Account) TableName() string {
	return "accounts"
}
