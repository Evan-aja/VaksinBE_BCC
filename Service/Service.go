package Service

import (
	"VaksinBE_BCC/Auth"
	"VaksinBE_BCC/User"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/service")
	r.GET("/daftar/swab", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying",
				"error":   err.Error.Error(),
			})
			return
		}
		var swab []Swab
		if err := db.Preload("Schedule").Find(&swab); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for services",
				"error":   err.Error.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "user data query successful",
			"data": gin.H{
				"name":          user.Name,
				"email":         user.Email,
				"handphone":     user.Handphone,
				"tanggal_lahir": user.TanggalLahir,
				"nik":           user.NIK,
				"nim":           user.NIM,
				"gender":        user.Gender,
			},
		})
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "System query sucessfull",
			"data":    swab,
		})
	})
	r.POST("/add/schedule", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying",
				"error":   err.Error.Error(),
			})
			return
		}
		if user.Email != "admin1@gmail.com" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sorry, this page is for administrator only",
			})
			return
		}
		sch := Schedule{}
		if err := c.BindJSON(&sch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Wrong data type might be inserted.",
				"error":   err.Error(),
			})
			return
		}
		if err := db.Create(&sch); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong when adding data to database",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Schedule added succesfully",
			"data":    sch,
		})
	})
	r.GET("add/swab", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying",
				"error":   err.Error.Error(),
			})
			return
		}
		if user.Email != "admin1@gmail.com" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sorry, this page is for administrator only",
			})
			return
		}
		var sch []Schedule
		if err := db.Find(&sch); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for schedules",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"message":   "ready to receive Swab Service addition. pick schedules by id",
			"schedules": sch,
		})
	})
	r.POST("/add/swab", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying",
				"error":   err.Error.Error(),
			})
			return
		}
		if user.Email != "admin1@gmail.com" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Sorry, this page is for administrator only",
			})
			return
		}
		sa := AddSwab{}
		if err := c.BindJSON(&sa); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please check the data you've submitted for swab",
			})
			return
		}
		sw := Swab{
			Type:   sa.Type,
			Cost:   sa.Cost,
			AdCost: sa.AdCost,
		}
		if err := db.Create(&sw); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please check your input data",
				"error":   err.Error.Error(),
			})
		}
		for _, val := range sa.ScheduleID {
			if err := db.Exec("INSERT INTO swab_schedule (swab_id,schedule_id) VALUES(?,?);", sw.ID, val); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "something went wrong in the database",
					"error":   err.Error.Error(),
				})
			}
		}
		result := Swab{}
		if err := db.Where("id=?", sw.ID).Preload("Schedule").Take(&result); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "error on query",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Service added",
			"data":    result,
		})
	})
}
