package User

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) {
	r.POST("/user/register", func(c *gin.Context) {
		var input UserRegister
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		hash := sha512.New()
		hash.Write([]byte(input.Password))
		pass := hex.EncodeToString(hash.Sum(nil))
		regist := User{
			Name:     input.Name,
			Email:    input.Email,
			Password: pass,
			Username: input.Username,
			NIK:      input.NIK,
			NIM:      input.NIM,
		}
		if err := db.Create(&regist); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Account created successfully",
			"error":   nil,
			"data": gin.H{
				"username": regist.Username,
				"email":    regist.Email,
			},
		})
	})
}
