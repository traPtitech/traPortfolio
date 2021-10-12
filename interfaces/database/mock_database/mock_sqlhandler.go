package mock_database //nolint:revive

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// TODO 適切な場所に書く
func init() {
	gorm.ErrRecordNotFound = repository.ErrNotFound
}

type MockSQLHandler struct {
	Conn *gorm.DB
	Mock sqlmock.Sqlmock
}

func NewMockSQLHandler(isValidDB bool) database.SQLHandler {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{SkipInitializeWithVersion: true}
	if isValidDB {
		conf.Conn = db
	}

	engine, err := gorm.Open(
		mysql.New(conf),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
	)
	if err != nil && isValidDB {
		log.Fatal(err)
	}

	return &MockSQLHandler{engine, mock}
}

func (handler *MockSQLHandler) Find(out interface{}, where ...interface{}) database.SQLHandler {
	db := handler.Conn.Find(out, where...)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) First(out interface{}, where ...interface{}) database.SQLHandler {
	db := handler.Conn.First(out, where...)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Create(value interface{}) database.SQLHandler {
	db := handler.Conn.Create(value)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Delete(value interface{}, where ...interface{}) database.SQLHandler {
	db := handler.Conn.Delete(value, where...)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	db := handler.Conn.Where(query, args...)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Model(value interface{}) database.SQLHandler {
	db := handler.Conn.Model(value)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Updates(values interface{}) database.SQLHandler {
	db := handler.Conn.Updates(values)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Begin() database.SQLHandler {
	tx := handler.Conn.Begin()
	return &MockSQLHandler{Conn: tx}
}

func (handler *MockSQLHandler) Commit() database.SQLHandler {
	db := handler.Conn.Commit()
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Preload(query string, args ...interface{}) database.SQLHandler {
	db := handler.Conn.Preload(query, args...)
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Rollback() database.SQLHandler {
	db := handler.Conn.Rollback()
	return &MockSQLHandler{Conn: db}
}

func (handler *MockSQLHandler) Transaction(fc func(database.SQLHandler) error) error {
	ffc := func(tx *gorm.DB) error {
		driver := &MockSQLHandler{Conn: tx}
		return fc(driver)
	}
	return handler.Conn.Transaction(ffc)
}

func (handler *MockSQLHandler) Error() error {
	return handler.Conn.Error
}

// Interface guards
var (
	_ database.SQLHandler = (*MockSQLHandler)(nil)
)
