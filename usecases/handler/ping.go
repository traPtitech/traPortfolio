package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (handler *PingHandler) Ping(c echo.Context) error {
	log.Println("ping received")
	return c.String(http.StatusOK, "pong")
}
