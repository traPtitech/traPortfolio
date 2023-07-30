package repository

import (
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/config"
)

const (
	ErrCodeInvalidConstraint = 1452
)

type dialector struct {
	gorm.Dialector
}

// override Translate(err error) error
func (d dialector) Translate(err error) error {
	if translater, ok := d.Dialector.(gorm.ErrorTranslator); ok {
		err = translater.Translate(err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.ErrNotFound
	}

	// 外部キー制約エラーの変換
	var mysqlErr *mysqldriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeInvalidConstraint {
		return repository.ErrInvalidArg
	}

	return err
}

func NewGormDB(conf *config.SQLConfig) (*gorm.DB, error) {
	d := dialector{mysql.New(mysql.Config{DSN: conf.Dsn()})}
	engine, err := gorm.Open(d, conf.GormConfig())
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
