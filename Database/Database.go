package database

import (
	"VaksinBE_BCC/User"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Open() *gorm.DB {
	var err error
	godotenv.Load(".env")
	db, err = gorm.Open(mysql.Open(os.Getenv("db")), &gorm.Config{})
	if err != nil {
		println(err.Error())
	}
	if err = db.AutoMigrate(&User.User{}, &User.UserPublic{}, &User.UserVacc{}); err != nil {
		fmt.Println(err.Error())
	}
	return db
}
