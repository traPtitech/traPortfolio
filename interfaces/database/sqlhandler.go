package database

import "gorm.io/gorm/clause"

type SQLHandler interface {
	Find(out interface{}, where ...interface{}) SQLHandler
	First(out interface{}, where ...interface{}) SQLHandler
	Create(value interface{}) SQLHandler
	Delete(value interface{}, where ...interface{}) SQLHandler
	Where(query interface{}, args ...interface{}) SQLHandler
	Model(value interface{}) SQLHandler
	Update(column string, value interface{}) SQLHandler
	Updates(values interface{}) SQLHandler
	Begin() SQLHandler
	Commit() SQLHandler
	Preload(query string, args ...interface{}) SQLHandler
	Rollback() SQLHandler
	Transaction(fc func(handler SQLHandler) error) error
	Clauses(conds ...clause.Expression) SQLHandler
	Ping() error

	Error() error
}
