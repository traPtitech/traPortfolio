package database

type SqlHandler interface {
	Exec(string, ...interface{}) SqlHandler
	Find(interface{}, ...interface{}) SqlHandler
	First(interface{}, ...interface{}) SqlHandler
	Raw(string, ...interface{}) SqlHandler
	Create(interface{}) SqlHandler
	Save(interface{}) SqlHandler
	Delete(interface{}) SqlHandler
	Where(interface{}, ...interface{}) SqlHandler
	Error() error
}
