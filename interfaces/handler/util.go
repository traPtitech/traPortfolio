package handler

import "time"

type duration struct {
	Since time.Time `json:"since"`
	Until time.Time `json:"until"`
}
