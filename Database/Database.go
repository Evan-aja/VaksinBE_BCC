package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Open() *gorm.DB {
	var err error
	db, err = gorm.Open(mysql.Open("root:123@tcp(127.0.0.1:41063)/bcc_backend?parseTime=true"), &gorm.Config{})
	if err != nil {
		println(err.Error())
	}
	return db
}
