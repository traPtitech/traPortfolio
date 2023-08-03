package infrastructure

import (
	"context"
	"errors"

	sqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/config"
)

type SQLHandler struct {
	conn *gorm.DB
}

func NewSQLHandler(db *gorm.DB) database.SQLHandler {
	return &SQLHandler{conn: db}
}

func NewGormDB(conf *config.SQLConfig) (*gorm.DB, error) {
	engine, err := gorm.Open(
		mysql.New(mysql.Config{DSNConfig: conf.DsnConfig()}),
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

func (h *SQLHandler) WithContext(ctx context.Context) database.SQLHandler {
	db := h.conn.WithContext(ctx)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Find(out interface{}, where ...interface{}) database.SQLHandler {
	db := h.conn.Find(out, where...)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) First(out interface{}, where ...interface{}) database.SQLHandler {
	db := h.conn.First(out, where...)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Create(value interface{}) database.SQLHandler {
	db := h.conn.Create(value)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Delete(value interface{}, where ...interface{}) database.SQLHandler {
	db := h.conn.Delete(value, where...)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	db := h.conn.Where(query, args...)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Model(value interface{}) database.SQLHandler {
	db := h.conn.Model(value)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Update(column string, value interface{}) database.SQLHandler {
	db := h.conn.Update(column, value)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Updates(values interface{}) database.SQLHandler {
	db := h.conn.Updates(values)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Begin() database.SQLHandler {
	tx := h.conn.Begin()
	return &SQLHandler{conn: tx}
}

func (h *SQLHandler) Commit() database.SQLHandler {
	db := h.conn.Commit()
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Preload(query string, args ...interface{}) database.SQLHandler {
	db := h.conn.Preload(query, args...)
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Rollback() database.SQLHandler {
	db := h.conn.Rollback()
	return &SQLHandler{conn: db}
}

func (h *SQLHandler) Transaction(fc func(database.SQLHandler) error) error {
	ffc := func(tx *gorm.DB) error {
		driver := &SQLHandler{conn: tx}
		return fc(driver)
	}
	return h.conn.Transaction(ffc)
}

func (h *SQLHandler) Ping() error {
	db, err := h.conn.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (h *SQLHandler) Limit(limit int) database.SQLHandler {
	db := h.conn.Limit(limit)
	return &SQLHandler{conn: db}
}

const (
	ErrCodeInvalidConstraint = 1452
)

func (h *SQLHandler) Error() error {
	err := h.conn.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.ErrNotFound
	}

	// 外部キー制約エラーの変換
	var mysqlErr *sqldriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeInvalidConstraint {
		return repository.ErrInvalidArg
	}

	return err
}

// Interface guards
var (
	_ database.SQLHandler = (*SQLHandler)(nil)
)
