package handler

import "time"

type duration struct {
	since time.Time `json:"since"`
	until time.Time `json:"until"`
}
