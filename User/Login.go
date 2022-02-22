package User

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) {
	r.POST("/user/login", func(c *gin.Context) {
		var input UserLogin
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		login := User{}
		if err := db.Where("email=?", input.Email).Take(&login); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "email does not exist",
				"error":   err.Error.Error(),
			})
			return
		}
		hash := sha512.New()
		hash.Write([]byte(input.Password))
		pass := hex.EncodeToString(hash.Sum(nil))
		if login.Password == pass {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})
			strToken, err := token.SignedString([]byte("GeneratorTokenSuperKompleks"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Something went wrong",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"succ":    true,
				"message": "Welcome, here's your token. don't lose it ;)",
				"data": gin.H{
					"email":    login.Email,
					"username": login.Username,
					"token":    strToken,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"succ":    false,
				"message": "Did you forget your own password?",
			})
			return
		}
	})
}
