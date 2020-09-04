package repository

type PingRepository interface {
	Ping() error
}
