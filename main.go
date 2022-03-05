package main

import (
	"VaksinBE_BCC/Dashboard"
	"VaksinBE_BCC/Database"
	"VaksinBE_BCC/Service"
	"VaksinBE_BCC/User"
	"VaksinBE_BCC/Vaccine"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

// main runner
func main() {
	fmt.Println("HI")
	db := database.Open()
	fmt.Println("Database terinisialisasi")

	r = gin.Default()
	r.Use(cors.Default())

	Dashboard.Routes(db, r)
	User.Routes(db, r)
	Vaccine.Routes(db, r)
	Service.Routes(db, r)

	fmt.Println("Router siap")
	fmt.Println("Server berjalan")
	if err := r.Run(); err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}
}
