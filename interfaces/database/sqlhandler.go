package database

type SQLHandler interface {
	Exec(sql string, values ...interface{}) SQLHandler
	Find(out interface{}, where ...interface{}) SQLHandler
	First(out interface{}, where ...interface{}) SQLHandler
	Raw(sql string, values ...interface{}) SQLHandler
	Create(value interface{}) SQLHandler
	Save(value interface{}) SQLHandler
	Delete(value interface{}, where ...interface{}) SQLHandler
	Where(query interface{}, args ...interface{}) SQLHandler
	Model(value interface{}) SQLHandler
	Updates(values interface{}) SQLHandler
	Begin() SQLHandler
	Commit() SQLHandler
	Joins(query string, args ...interface{}) SQLHandler
	Scan(dest interface{}) SQLHandler
	Select(query interface{}, args ...interface{}) SQLHandler
	Rollback() SQLHandler
	Transaction(fc func(handler SQLHandler) error) error

	Error() error
}
