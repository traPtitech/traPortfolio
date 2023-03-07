package mock_database //nolint:revive

import (
	"context"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type MockSQLHandler struct {
	Conn *gorm.DB
	Mock sqlmock.Sqlmock
}

func NewMockSQLHandler() *MockSQLHandler {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{SkipInitializeWithVersion: true}
	conf.Conn = db

	engine, err := gorm.Open(mysql.New(conf), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &MockSQLHandler{engine, mock}
}

func (h *MockSQLHandler) WithContext(ctx context.Context) database.SQLHandler {
	db := h.Conn.WithContext(ctx)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Find(out interface{}, where ...interface{}) database.SQLHandler {
	db := h.Conn.Find(out, where...)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) First(out interface{}, where ...interface{}) database.SQLHandler {
	db := h.Conn.First(out, where...)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Create(value interface{}) database.SQLHandler {
	db := h.Conn.Create(value)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Delete(value interface{}, where ...interface{}) database.SQLHandler {
	db := h.Conn.Delete(value, where...)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Update(column string, value interface{}) database.SQLHandler {
	db := h.Conn.Update(column, value)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Where(query interface{}, args ...interface{}) database.SQLHandler {
	db := h.Conn.Where(query, args...)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Model(value interface{}) database.SQLHandler {
	db := h.Conn.Model(value)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Updates(values interface{}) database.SQLHandler {
	db := h.Conn.Updates(values)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Begin() database.SQLHandler {
	tx := h.Conn.Begin()
	return &MockSQLHandler{Conn: tx}
}

func (h *MockSQLHandler) Commit() database.SQLHandler {
	db := h.Conn.Commit()
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Preload(query string, args ...interface{}) database.SQLHandler {
	db := h.Conn.Preload(query, args...)
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Rollback() database.SQLHandler {
	db := h.Conn.Rollback()
	return &MockSQLHandler{Conn: db}
}

func (h *MockSQLHandler) Transaction(fc func(database.SQLHandler) error) error {
	ffc := func(tx *gorm.DB) error {
		driver := &MockSQLHandler{Conn: tx}
		return fc(driver)
	}
	return h.Conn.Transaction(ffc)
}

func (h *MockSQLHandler) Ping() error {
	db, err := h.Conn.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (h *MockSQLHandler) Error() error {
	return h.Conn.Error
}

func (h *MockSQLHandler) Limit(limit int) database.SQLHandler {
	db := h.Conn.Limit(limit)
	return &MockSQLHandler{Conn: db}
}

// Interface guards
var (
	_ database.SQLHandler = (*MockSQLHandler)(nil)
)
