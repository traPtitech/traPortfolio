package infrastructure

import (
	"fmt"
	"os"
	"strconv"

	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type SqlHandler struct {
	Conn *gorm.DB
}

func NewSqlHandler() database.SqlHandler {
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
		host = "localhost"
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
	}

	db = db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Set("gorm:table_options", "CHARSET=utf8mb4")

	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(16)
	db.BlockGlobalUpdate(true)

	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = db
	return sqlHandler
}

func (handler *SqlHandler) Find(out interface{}, where ...interface{}) database.SqlHandler {
	res := handler.Conn.Find(out, where...)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Exec(sql string, values ...interface{}) database.SqlHandler {
	res := handler.Conn.Exec(sql, values...)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) First(out interface{}, where ...interface{}) database.SqlHandler {
	res := handler.Conn.First(out, where...)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Raw(sql string, values ...interface{}) database.SqlHandler {
	res := handler.Conn.Raw(sql, values...)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Create(value interface{}) database.SqlHandler {
	res := handler.Conn.Create(value)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Save(value interface{}) database.SqlHandler {
	res := handler.Conn.Save(value)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Delete(value interface{}) database.SqlHandler {
	res := handler.Conn.Delete(value)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Where(query interface{}, args ...interface{}) database.SqlHandler {
	res := handler.Conn.Where(query, args...)
	handler.Conn = res
	return handler
}

func (handler *SqlHandler) Error() error {
	return handler.Conn.Error
}
