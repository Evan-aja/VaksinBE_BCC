package User

import (
	auth "VaksinBE_BCC/Auth"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/user")
	r.GET("/", auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if result := db.Where("id=?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"succ":    false,
				"message": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"succ":    true,
			"message": "success",
			"data": gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"username": user.Username,
				"nik":      user.NIK,
				"nim":      user.NIM,
			},
		})
	})
	r.POST("/register", func(c *gin.Context) {
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
		registPublic := UserPublic{
			Name:     input.Name,
			Email:    input.Email,
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
		if err := db.Create(&registPublic); err.Error != nil {
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
	r.POST("/login", func(c *gin.Context) {
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
	r.GET("/search", func(c *gin.Context) {
		name, nameExists := c.GetQuery("name")
		email, emailExists := c.GetQuery("email")
		username, usernameExists := c.GetQuery("username")
		if !nameExists && !usernameExists && !emailExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Data incorrect.",
			})
			return
		}
		var query []UserPublic
		dbtmp := db
		if nameExists {
			dbtmp = dbtmp.Where("name LIKE ?", "%"+name+"%")
		}
		if emailExists {
			dbtmp = dbtmp.Where("email LIKE ?", "%"+email+"%")
		}
		if usernameExists {
			dbtmp = dbtmp.Where("username LIKE ?", "%"+username+"%")
		}
		if result := dbtmp.Find(&query); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something wrong happened",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"succ": true,
			"data": gin.H{
				"query": gin.H{
					"name":     name,
					"email":    email,
					"username": username,
				},
				"result": query,
			},
		})
	})
	r.DELETE("/", auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		userpub := UserPublic{}
		if result := db.Where("id=?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"succ":    false,
				"message": result.Error.Error(),
			})
			return
		}
		if result := db.Where("id=?", id).Take(&userpub); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"succ":    false,
				"message": result.Error.Error(),
			})
			return
		}
		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"succ":    false,
				"message": result.Error.Error(),
			})
			return
		}
		if result := db.Delete(&userpub); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"succ":    false,
				"message": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"succ":    true,
			"deleted": userpub,
		})
	})

}
