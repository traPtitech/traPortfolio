package repository

import "log"

type PingRepository struct {
}

func NewPingRepository() *PingRepository {
	return &PingRepository{}
}

func (repo *PingRepository) Ping() (err error) {
	log.Println("pinged!!")
	return
}
