package database

type SQLHandler interface {
	Exec(string, ...interface{}) SQLHandler
	Find(interface{}, ...interface{}) SQLHandler
	First(interface{}, ...interface{}) SQLHandler
	Raw(string, ...interface{}) SQLHandler
	Create(interface{}) SQLHandler
	Save(interface{}) SQLHandler
	Delete(interface{}) SQLHandler
	Where(interface{}, ...interface{}) SQLHandler
	Error() error
}
