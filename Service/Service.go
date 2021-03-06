package Service

import (
	"VaksinBE_BCC/Auth"
	"VaksinBE_BCC/User"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/service")
	// Untuk mendapatkan data pengguna dan daftar jadwal swab yang tersedia
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
			"success":            true,
			"message":            "System query sucessfull",
			"available_schedule": swab,
		})
	})
	// Untuk mendaftarkan pengguna menggunakan data pengguna dan daftar jadwal swab yang dipilih
	r.POST("/daftar/swab", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying user",
				"error":   err.Error.Error(),
			})
			return
		}
		transac := TransactionSwab{}
		if err := c.BindJSON(&transac); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong while binding",
				"error":   err.Error(),
			})
			return
		}
		if user.Email != transac.Email || user.Gender != transac.Gender || user.Handphone != transac.Handphone || user.NIK != transac.NIK || user.NIM != transac.NIM || user.Name != transac.Name {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "user data and transaction data do not line up. please update your data on profile page.",
			})
			return
		}
		if transac.NIM != "" {
			transac.AdCost = 0
		}
		if transac.Date.Unix() > time.Now().Add(time.Hour*24*7).Unix() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Booking time must not exceed 7 days",
			})
			return
		} else if transac.Date.Unix() < time.Now().Unix() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Booking time must not be on past date",
			})
			return
		}

		transac.Paid = false
		transac.UserID = user.ID
		if err := db.Create(&transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Transaction cannot be processed. please contact support for more help",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Where("id=?", transac.ID).Preload("Schedule").Take(&transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server failed to query your appointment",
				"error":   err.Error.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction has been completed successfully. please fulfill the payment immediately after",
			"data":    transac,
		})
	})

	// Untuk mendapatkan riwayat pengguna dalam penggunaan layanan untuk swab test
	r.GET("/riwayat/swab", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var transac []TransactionSwab
		if err := db.Preload("Schedule").Find(&transac).Where("user_id=?", id); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for history",
				"error":   err.Error.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "System query sucessfull",
			"history": transac,
		})
	})

	// Untuk mengubah data swab test yang telah dipesan pengguna (untuk bukti pembayaran)
	r.PATCH("/update/swab", Auth.Authorization(), func(c *gin.Context) {
		// id, _ := c.Get("id")
		transac := TransactionSwab{}
		transup := TransSwabUp{}
		if err := c.BindJSON(&transup); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		if err := db.Preload("Schedule").Find(&transac).Where("id=?", transup.ID); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for history",
				"error":   err.Error.Error(),
			})
			return
		}
		transac = TransactionSwab{
			ID:           transac.ID,
			UserID:       transac.UserID,
			User:         transac.User,
			Name:         transac.Name,
			Email:        transac.Email,
			Handphone:    transac.Handphone,
			TanggalLahir: transac.TanggalLahir,
			NIK:          transac.NIK,
			NIM:          transac.NIM,
			Gender:       transac.Gender,
			Type:         transac.Type,
			Date:         transac.Date,
			Cost:         transac.Cost,
			AdCost:       transac.AdCost,
			IDSchedule:   transac.IDSchedule,
			Schedule:     transac.Schedule,
			Paid:         transup.Paid,
			Receipt:      transup.Receipt,
		}
		if err := db.Model(&transac).Where("id=?", transac.ID).Updates(transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database or id is incorrect",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "System query sucessfull",
			"history": transac,
		})
	})

	// Untuk mendapatkan data pengguna dan daftar jadwal poliklinik yang tersedia
	r.GET("/daftar/poly", Auth.Authorization(), func(c *gin.Context) {
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
		var poly []Poly
		if err := db.Preload("Schedule").Find(&poly); err.Error != nil {
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
			"success":            true,
			"message":            "System query sucessfull",
			"available_schedule": poly,
		})
	})
	// Untuk mendaftarkan pengguna menggunakan data pengguna dan daftar jadwal poliklinik yang dipilih
	r.POST("/daftar/poly", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User.User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying user",
				"error":   err.Error.Error(),
			})
			return
		}
		transac := TransactionPoly{}
		if err := c.BindJSON(&transac); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong while binding",
				"error":   err.Error(),
			})
			return
		}
		if user.Email != transac.Email || user.Gender != transac.Gender || user.Handphone != transac.Handphone && user.NIK != transac.NIK || user.NIM != transac.NIM || user.Name != transac.Name {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "user data and transaction data does not line up. please update your data on profile page.",
			})
			return
		}
		if transac.NIM != "" {
			transac.AdCost = 0
		}
		if transac.Date.Unix() > time.Now().Add(time.Hour*24*7).Unix() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Booking time must not exceed 7 days",
			})
			return
		} else if transac.Date.Unix() < time.Now().Unix() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Booking time must not be on past date",
			})
			return
		}

		transac.Paid = false
		transac.UserID = user.ID
		if err := db.Create(&transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Transaction cannot be processed. please contact support for more help",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Where("id=?", transac.ID).Preload("Schedule").Take(&transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server failed to query your appointment",
				"error":   err.Error.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction has been completed successfully. please fulfill the payment immediately after",
			"data":    transac,
		})
	})

	// Untuk mendapatkan riwayat pengguna dalam penggunaan layanan untuk poliklinik
	r.GET("/riwayat/poly", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var transac []TransactionPoly
		if err := db.Preload("Schedule").Find(&transac).Where("user_id=?", id); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for history",
				"error":   err.Error.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "System query sucessfull",
			"history": transac,
		})
	})

	// Untuk mengubah data layanan poliklinik yang telah dipesan pengguna (untuk bukti pembayaran)
	r.PATCH("/update/poly", Auth.Authorization(), func(c *gin.Context) {
		transac := TransactionPoly{}
		transup := TransPolyUp{}
		if err := c.BindJSON(&transup); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		if err := db.Preload("Schedule").Find(&transac).Where("id=?", transup.ID); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong while querying for history",
				"error":   err.Error.Error(),
			})
			return
		}
		transac = TransactionPoly{
			ID:           transac.ID,
			UserID:       transac.UserID,
			User:         transac.User,
			Name:         transac.Name,
			Email:        transac.Email,
			Handphone:    transac.Handphone,
			TanggalLahir: transac.TanggalLahir,
			NIK:          transac.NIK,
			NIM:          transac.NIM,
			Gender:       transac.Gender,
			Type:         transac.Type,
			Date:         transac.Date,
			Cost:         transac.Cost,
			AdCost:       transac.AdCost,
			IDSchedule:   transac.IDSchedule,
			Schedule:     transac.Schedule,
			Description:  transac.Description,
			Paid:         transup.Paid,
			Receipt:      transup.Receipt,
		}
		if err := db.Model(&transac).Where("id=?", transac.ID).Updates(transac); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database or id is incorrect",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "System query sucessfull",
			"history": transac,
		})
	})

	// menambahkan jadwal yang tersedia. dalam string jam (08:00, 09:00 dsb)
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

	// mendapatkan jadwal yang tersedia untuk penambahan layanan swab test
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
	// menambahkan layanan swab test berdasarkan data yang diberikan dan jadwal yang tersedia
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

	// mendapatkan jadwal yang tersedia untuk penambahan layanan poliklinik
	r.GET("add/poly", Auth.Authorization(), func(c *gin.Context) {
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
			"message":   "ready to receive Poly Service addition. pick schedules by id",
			"schedules": sch,
		})
	})
	// menambahkan layanan poliklinik berdasarkan data yang diberikan dan jadwal yang tersedia
	r.POST("/add/poly", Auth.Authorization(), func(c *gin.Context) {
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
		pa := AddPoly{}
		if err := c.BindJSON(&pa); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please check the data you've submitted for swab",
			})
			return
		}
		pw := Poly{
			Type:        pa.Type,
			Cost:        pa.Cost,
			AdCost:      pa.AdCost,
			Description: pa.Description,
		}
		if err := db.Create(&pw); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please check your input data",
				"error":   err.Error.Error(),
			})
		}
		for _, val := range pa.ScheduleID {
			if err := db.Exec("INSERT INTO poly_schedule (poly_id,schedule_id) VALUES(?,?);", pw.ID, val); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "something went wrong in the database",
					"error":   err.Error.Error(),
				})
			}
		}
		result := Poly{}
		if err := db.Where("id=?", pw.ID).Preload("Schedule").Take(&result); err.Error != nil {
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
