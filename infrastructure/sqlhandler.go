package infrastructure

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type SQLConfig struct {
	user     string
	password string
	host     string
	port     int
	dbname   string
}

func NewSQLConfig(user, password, host, dbname string, port int) SQLConfig {
	return SQLConfig{
		user,
		password,
		host,
		port,
		dbname,
	}
}

type SQLHandler struct {
	conn *gorm.DB
}

func NewSQLHandler(conf *SQLConfig) (database.SQLHandler, error) {
	engine, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&collation=utf8mb4_general_ci&loc=Local", conf.user, conf.password, conf.host, conf.port, conf.dbname),
	}), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Truncate(time.Microsecond)
		},
	})
	if err != nil {
		// return fmt.Errorf("failed to connect database: %v", err)s
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

	sqlHandler := new(SQLHandler)
	sqlHandler.conn = engine
	return sqlHandler, nil
}

// initDB データベースのスキーマを更新
func initDB(db *gorm.DB) error {
	// gormのエラーの上書き
	gorm.ErrRecordNotFound = repository.ErrNotFound
	// db.Logger = db.Logger.LogMode(logger.Info)
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

func (handler *SQLHandler) Update(column string, values interface{}) database.SQLHandler {
	db := handler.conn.Update(column, values)
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

func (handler *SQLHandler) Error() error {
	return handler.conn.Error
}

// Interface guards
var (
	_ database.SQLHandler = (*SQLHandler)(nil)
)
