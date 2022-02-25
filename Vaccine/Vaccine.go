package Vaccine

import (
	"VaksinBE_BCC/Auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/vaccine")
	r.PATCH("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var input Vaccine
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		vaksin := Vaccine{}
		if err := db.Where("id=?", id).Take(&vaksin); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if input.Dosis1 && !input.Dosis2 && !input.Booster {
			input.Dosis2 = false
			input.Booster = false
		}
		if input.Dosis2 && !input.Dosis1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid, Dosis 2 hanya bisa diambil jika sudah mengambil Dosis 1",
			})
			return
		}
		if input.Dosis2 && !input.Booster {
			input.Dosis1 = true
			input.Booster = false
		}
		if input.Booster && !input.Dosis2 || input.Booster && !input.Dosis1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid, Dosis Booster hanya bisa diambil jika sudah mengambil Dosis 1 dan 2",
			})
			return
		}
		if input.Booster {
			input.Dosis1 = true
			input.Dosis2 = true
		}
		vaksin = Vaccine{
			ID:      vaksin.ID,
			Dosis1:  input.Dosis1,
			Dosis2:  input.Dosis2,
			Booster: input.Booster,
		}
		err := db.Model(&vaksin).Updates(map[string]interface{}{"id": vaksin.ID, "dosis1": vaksin.Dosis1, "dosis2": vaksin.Dosis2, "booster": vaksin.Booster})
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   err.Error.Error(),
			})
			return
		}
		if err.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No data has been updated (could also be if the same data as before was sent)",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    vaksin,
			"input":   input,
		})
	})
}
