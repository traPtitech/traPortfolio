package infrastructure

import (
	"encoding/json"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
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

func (handler *SQLHandler) Find(out interface{}, where ...interface{}) database.SQLHandler {
	db := handler.conn.Find(out, where...)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) First(out interface{}, where ...interface{}) database.SQLHandler {
	db := handler.conn.First(out, where...)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Create(value interface{}) database.SQLHandler {
	db := handler.conn.Create(value)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Delete(value interface{}, where ...interface{}) database.SQLHandler {
	db := handler.conn.Delete(value, where...)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	db := handler.conn.Where(query, args...)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Model(value interface{}) database.SQLHandler {
	db := handler.conn.Model(value)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Update(column string, value interface{}) database.SQLHandler {
	db := handler.conn.Update(column, value)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Updates(values interface{}) database.SQLHandler {
	db := handler.conn.Updates(values)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Begin() database.SQLHandler {
	tx := handler.conn.Begin()
	return &SQLHandler{conn: tx}
}

func (handler *SQLHandler) Commit() database.SQLHandler {
	db := handler.conn.Commit()
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Preload(query string, args ...interface{}) database.SQLHandler {
	db := handler.conn.Preload(query, args...)
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Rollback() database.SQLHandler {
	db := handler.conn.Rollback()
	return &SQLHandler{conn: db}
}

func (handler *SQLHandler) Transaction(fc func(database.SQLHandler) error) error {
	ffc := func(tx *gorm.DB) error {
		driver := &SQLHandler{conn: tx}
		return fc(driver)
	}
	return handler.conn.Transaction(ffc)
}

func (handler *SQLHandler) Ping() error {
	db, err := handler.conn.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

type GormErr struct {
	Number int `json:"Number"`
}

func (handler *SQLHandler) Error() error {
	err := handler.conn.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.ErrNoRows
	}

	// 外部キー制約エラーの変換
	b, _ := json.Marshal(err)
	var newError GormErr
	error := json.Unmarshal((b), &newError)
	if error != nil {
		return err
	}
	if newError.Number == 1452 {
		return database.ErrInvalidArgument
	}

	return err
}

// Interface guards
var (
	_ database.SQLHandler = (*SQLHandler)(nil)
)
