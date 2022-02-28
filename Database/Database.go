package database

import (
	"VaksinBE_BCC/User"
	"VaksinBE_BCC/Vaccine"
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
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))),
		&gorm.Config{})
	if err != nil {
		println(err.Error())
	}
	if err = db.AutoMigrate(&User.User{}, &User.UserPublic{}, &Vaccine.Vaccine{}, &Vaccine.VaccProof{}); err != nil {
		fmt.Println(err.Error())
	}
	return db
}
