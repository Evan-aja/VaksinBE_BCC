package main

import (
	"VaksinBE_BCC/User"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:123@tcp(127.0.0.1:41063)/bcc_backend?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	if err = db.AutoMigrate(&User.User{}); err != nil {
		return err
	}
	return nil
}

func initRouter() {
	User.Register(db)

}

func main() {
	fmt.Println("HI")
	if err := InitDB(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Database terinisialisasi")
}
