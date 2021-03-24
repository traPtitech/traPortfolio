package domain

import "time"

type Duration struct {
	Since time.Time `gorm:"precision:6"`
	Until time.Time `gorm:"precision:6"`
}
