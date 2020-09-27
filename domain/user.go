package domain

import "time"

type User struct {
	ID          uint      `gorm:"type:char(36);not null;primary_key"`
	Description string    `sql:"type:TEXT COLLATE utf8mb4_bin NOT NULL"`
	Check       bool      `gorm:"type:boolean;not null;default:false"`
	Name        string    `gorm:"type:varchar(32);not null;unique"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}
