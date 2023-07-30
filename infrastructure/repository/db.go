package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/util/config"
)

func NewGormDB(conf *config.SQLConfig) (*gorm.DB, error) {
	engine, err := gorm.Open(
		mysql.New(mysql.Config{DSN: conf.Dsn()}),
		conf.GormConfig(),
	)
	if err != nil {
		return nil, err
	}

	db, err := engine.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(16)
	if err := initDB(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

// initDB データベースのスキーマを更新
func initDB(db *gorm.DB) error {
	_, err := migration.Migrate(db, migration.AllTables())
	if err != nil {
		return err
	}
	return nil
}

const (
	ErrCodeInvalidConstraint = 1452
)

// func (h *SQLHandler) Error() error {
// 	err := h.conn.Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return repository.ErrNotFound
// 	}

// 	// 外部キー制約エラーの変換
// 	var mysqlErr *sqldriver.MySQLError
// 	if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeInvalidConstraint {
// 		return repository.ErrInvalidArg
// 	}

// 	return err
// }
