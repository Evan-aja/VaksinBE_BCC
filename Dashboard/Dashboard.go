package Dashboard

import (
	"VaksinBE_BCC/Vaccine"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/dashboard")
	r.GET("/", func(c *gin.Context) {
		var id = int64(0)
		var dosis1 = int64(0)
		var dosis2 = int64(0)
		var booster = int64(0)
		if err := db.Model(&Vaccine.Vaccine{}).Count(&id); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong on the database",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Model(&Vaccine.Vaccine{}).Where("dosis1=?", 1).Count(&dosis1); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong on the database",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Model(&Vaccine.Vaccine{}).Where("dosis2=?", 1).Count(&dosis2); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong on the database",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Model(&Vaccine.Vaccine{}).Where("booster=?", 1).Count(&booster); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong on the database",
				"error":   err.Error.Error(),
			})
			return
		}
		persen1 := fmt.Sprintf("%.2f", float64(dosis1*100/id))
		persen2 := fmt.Sprintf("%.2f", float64(dosis2*100/id))
		persen3 := fmt.Sprintf("%.2f", float64(booster*100/id))
		c.JSON(http.StatusOK, gin.H{
			"success":            true,
			"jumlah pengguna":    id,
			"jumlah dosis1":      dosis1,
			"persentase dosis1":  persen1,
			"jumlah dosis2":      dosis2,
			"persentase dosis2":  persen2,
			"jumlah booster":     booster,
			"persentase booster": persen3,
		})
	})
}
