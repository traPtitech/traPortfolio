package handler

import (
	"time"

	"github.com/traPtitech/traPortfolio/util/optional"
)

type Duration struct {
	Since time.Time `json:"since"`
	Until time.Time `json:"until"`
}

type OptionalDuration struct {
	Since optional.Time `json:"since"`
	Until optional.Time `json:"until"`
}
