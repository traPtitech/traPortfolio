package infrastructure

import (
	"fmt"
	"os"
	"strconv"

	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type SQLHandler struct {
	Conn *gorm.DB
}

func NewSQLHandler() (database.SQLHandler, error) {
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
		Set("gorm:table_options", "CHARSET=utf8mb4")

	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(16)
	db.BlockGlobalUpdate(true)

	sqlHandler := new(SQLHandler)
	sqlHandler.Conn = db
	return sqlHandler, nil
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

func (handler *SQLHandler) Delete(value interface{}) database.SQLHandler {
	res := handler.Conn.Delete(value)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	res := handler.Conn.Where(query, args...)
	handler.Conn = res
	return handler
}

func (handler *SQLHandler) Error() error {
	return handler.Conn.Error
}
