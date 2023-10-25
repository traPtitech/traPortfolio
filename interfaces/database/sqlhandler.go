package database

import "context"

type SQLHandler interface {
	WithContext(ctx context.Context) SQLHandler
	Find(out interface{}, where ...interface{}) SQLHandler
	First(out interface{}, where ...interface{}) SQLHandler
	Create(value interface{}) SQLHandler
	Delete(value interface{}, where ...interface{}) SQLHandler
	Where(query interface{}, args ...interface{}) SQLHandler
	Select(query interface{}, args ...interface{}) SQLHandler
	Model(value interface{}) SQLHandler
	Update(column string, value interface{}) SQLHandler
	Updates(values interface{}) SQLHandler
	Begin() SQLHandler
	Commit() SQLHandler
	Preload(query string, args ...interface{}) SQLHandler
	Rollback() SQLHandler
	Transaction(fc func(handler SQLHandler) error) error
	Ping() error
	Limit(limit int) SQLHandler
	Error() error
}
