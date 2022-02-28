package Vaccine

import (
	"VaksinBE_BCC/Auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/vaccine")

	// vaccine updates, needs all value to be filled accordingly. details in the code
	r.PATCH("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var input VaccineIn
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		vaksin := Vaccine{}
		bukti := VaccProof{}
		if err := db.Where("id=?", id).Take(&vaksin); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		// if dosis1 is true and everything else is false, set everything to null or false
		if input.Dosis1 && !input.Dosis2 && !input.Booster {
			// if proof is null then bad request
			if input.Bukti1 == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "No link for proof has been sent",
				})
				return
			}
			input.Dosis2 = false
			input.Bukti2 = ""
			input.Booster = false
			input.Bukti3 = ""
		}
		// if dosis2 is true and dosis1 is false, returns bad request
		if input.Dosis2 && !input.Dosis1 {
			input.Bukti1 = ""
			input.Bukti2 = ""
			input.Bukti3 = ""
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid, Dosis 2 hanya bisa diambil jika sudah mengambil Dosis 1",
			})
			return
		}
		// if dosis2 is true and booster is false, will set everything except booster to input value
		if input.Dosis2 && !input.Booster {
			// if proof is null then bad request
			if input.Bukti2 == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "No link for proof has been sent",
				})
				return
			}
			input.Dosis1 = true
			input.Booster = false
			input.Bukti3 = ""
		}
		// if biister is true but everything else is false, send bad request
		if input.Booster && !input.Dosis2 || input.Booster && !input.Dosis1 {
			input.Bukti1 = ""
			input.Bukti2 = ""
			input.Bukti3 = ""
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid, Dosis Booster hanya bisa diambil jika sudah mengambil Dosis 1 dan 2",
			})
			return
		}
		// if booster is true, and everything else did not violate previous rules, set to true
		if input.Booster {
			// if proof is null then bad request
			if input.Bukti3 == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "No link for proof has been sent",
				})
				return
			}
			input.Dosis1 = true
			input.Dosis2 = true
		}
		vaksin = Vaccine{
			ID:      vaksin.ID,
			Dosis1:  input.Dosis1,
			Dosis2:  input.Dosis2,
			Booster: input.Booster,
		}
		bukti = VaccProof{
			IDVaccine: vaksin.ID,
			Dosis1:    input.Bukti1,
			Dosis2:    input.Bukti2,
			Booster:   input.Bukti3,
		}
		err := db.Model(&vaksin).Updates(map[string]interface{}{"id": vaksin.ID, "dosis1": vaksin.Dosis1, "dosis2": vaksin.Dosis2, "booster": vaksin.Booster})
		ers := db.Model(&bukti).Where("id_vaccine=?", vaksin.ID).Updates(map[string]interface{}{"id_vaccine": vaksin.ID, "dosis1": bukti.Dosis1, "dosis2": bukti.Dosis2, "booster": bukti.Booster})
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the base database.",
				"error":   err.Error.Error(),
			})
			return
		} else if ers.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the proof database.",
				"error":   ers.Error.Error(),
			})
			return
		} else if err.RowsAffected < 1 || ers.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No data (value or links) has been updated (could also be if the same data as before was sent)",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data": gin.H{
				"id":           vaksin.ID,
				"dosis1":       vaksin.Dosis1,
				"dosis2":       vaksin.Dosis2,
				"booster":      vaksin.Booster,
				"buktidosis1":  bukti.Dosis1,
				"buktidosis2":  bukti.Dosis2,
				"buktibooster": bukti.Booster,
			},
			"input": input,
		})
	})
}
