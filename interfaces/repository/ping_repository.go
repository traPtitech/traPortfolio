package repository

import "log"

type PingRepository struct {
}

func (repo *PingRepository) Ping() (err error) {
	log.Println("pinged!!")
	return
}
