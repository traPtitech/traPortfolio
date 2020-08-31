package usecase

type PingRepository interface {
	Ping() error
}
