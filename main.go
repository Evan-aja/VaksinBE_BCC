package main

import (
	"VaksinBE_BCC/Database"
	"VaksinBE_BCC/User"
	"fmt"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	fmt.Println("HI")
	db := database.Open()
	fmt.Println("Database terinisialisasi")

	r = gin.Default()

	User.Routes(db, r)

	fmt.Println("Router siap")
	fmt.Println("Server berjalan")
	if err := r.Run(); err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}
}
