package infrastructure

import (
	"fmt"
	"os"
	"strconv"

	gorm "github.com/jinzhu/gorm"
	//
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type SQLHandler struct {
	Conn *gorm.DB
}

func NewSQLHandler() (*SQLHandler, error) {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "mysql"
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "portfolio"
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname))
	if err != nil {
		// return fmt.Errorf("failed to connect database: %v", err)s
		return nil, err
	}

	db = db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Set("gorm:table_options", "CHARSET=utf8mb4 COLLATE=utf8mb4_bin")

	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(16)
	db.BlockGlobalUpdate(true)
	if err := initDB(db); err != nil {
		return nil, err
	}

	sqlHandler := new(SQLHandler)
	sqlHandler.Conn = db
	return sqlHandler, nil
}

// initDB データベースのスキーマを更新
func initDB(db *gorm.DB) error {
	// gormのエラーの上書き
	gorm.ErrRecordNotFound = repository.ErrNotFound
	// db.LogMode(true)
	if err := migration.Migrate(db, migration.AllTables()); err != nil {
		return err
	}
	return nil
}

func (handler *SQLHandler) Connect(dialect string, args ...interface{}) error {
	db, err := gorm.Open(dialect, args...)
	if err != nil {
		return err
	}
	handler.Conn = db
	return nil
}

func (handler *SQLHandler) Find(out interface{}, where ...interface{}) database.SQLHandler {
	res := handler.Conn.Find(out, where...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Exec(sql string, values ...interface{}) database.SQLHandler {
	res := handler.Conn.Exec(sql, values...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) First(out interface{}, where ...interface{}) database.SQLHandler {
	res := handler.Conn.First(out, where...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Raw(sql string, values ...interface{}) database.SQLHandler {
	res := handler.Conn.Raw(sql, values...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Create(value interface{}) database.SQLHandler {
	res := handler.Conn.Create(value)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Save(value interface{}) database.SQLHandler {
	res := handler.Conn.Save(value)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Delete(value interface{}, where ...interface{}) database.SQLHandler {
	res := handler.Conn.Delete(value, where)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	res := handler.Conn.Where(query, args...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Model(value interface{}) database.SQLHandler {
	res := handler.Conn.Model(value)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Updates(values interface{}) database.SQLHandler {
	res := handler.Conn.Updates(values)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Begin() database.SQLHandler {
	tx := handler.Conn.Begin()
	return &SQLHandler{
		Conn: tx,
	}
}

func (handler *SQLHandler) Commit() database.SQLHandler {
	res := handler.Conn.Commit()
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Joins(query string, args ...interface{}) database.SQLHandler {
	res := handler.Conn.Joins(query, args)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Rollback() database.SQLHandler {
	res := handler.Conn.Rollback()
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Transaction(fc func(database.SQLHandler) error) error {
	ffc := func(tx *gorm.DB) error {
		driver := &SQLHandler{Conn: tx}
		return fc(driver)
	}
	return handler.Conn.Transaction(ffc)
}

func (handler *SQLHandler) Error() error {
	return handler.Conn.Error
}

// Interface guards
var (
	_ database.SQLHandler = (*SQLHandler)(nil)
)
