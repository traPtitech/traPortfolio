// Package migration migrate current struct
package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v1 unique_index:idx_room_uniqueの削除
func v1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",
		Migrate: func(db *gorm.DB) error {
			return db.
				Table("portfolio").
				Error
		},
	}
}
