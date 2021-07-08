package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrate execute migrations
// 初回実行でスキーマが初期化された場合、initでtrueを返します
func Migrate(db *gorm.DB, tables []interface{}) (init bool, err error) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, Migrations())

	m.InitSchema(func(db *gorm.DB) error {
		init = true

		return db.AutoMigrate(AllTables()...)
	})
	err = m.Migrate()
	return
}
