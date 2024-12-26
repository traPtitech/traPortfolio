// Package migration migrate current struct
package migration

import (
	"github.com/traPtitech/traPortfolio/internal/domain"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

// v3 UserにdisplayNameを追加
func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(db *gorm.DB) error {
			if err := db.AutoMigrate(&v3User{}); err != nil {
				return err
			}

			return db.
				Table("portfolio").
				Error
		},
	}
}

type v3User struct {
	ID          uuid.UUID        `gorm:"type:char(36);not null;primaryKey"`
	Description string           `gorm:"type:text;not null"`
	Check       bool             `gorm:"type:boolean;not null;default:false"`
	Name        string           `gorm:"type:varchar(32);not null;unique"`
	DisplayName *string          `gorm:"type:varchar(32)"` // 追加
	State       domain.TraQState `gorm:"type:tinyint(1);not null"`
	CreatedAt   time.Time        `gorm:"precision:6"`
	UpdatedAt   time.Time        `gorm:"precision:6"`

	Accounts []*model.Account `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*v3User) TableName() string {
	return "projects"
}
