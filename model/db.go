package model

import (
	"fmt"
	"os"
	"strconv"

	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func Setup() error {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "portfolio"
	}

	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname))
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	db = db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Set("gorm:table_options", "CHARSET=utf8mb4")

	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(16)
	db.BlockGlobalUpdate(true)

	return nil
}
