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
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "success",
			"data": gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"email":    user.Email,
				"username": user.Username,
				"nik":      user.NIK,
				"nim":      user.NIM,
			},
		})
	})
	r.GET("/search", func(c *gin.Context) {
		name, nameExists := c.GetQuery("name")
		email, emailExists := c.GetQuery("email")
		username, usernameExists := c.GetQuery("username")
		if !nameExists && !usernameExists && !emailExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong",
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
		if err := dbtmp.Find(&query); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
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
		regist := User{
			Name:     input.Name,
			Email:    input.Email,
			Password: hash(input.Password),
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
		registVacc := UserVacc{
			Dosis1:  false,
			Dosis2:  false,
			Booster: false,
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
		if err := db.Create(&registVacc); err.Error != nil {
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
			if err = db.Where("username=?", input.Username).Take(&login); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Email does not exist",
					"error":   err.Error.Error(),
				})
				return
			}
		}
		if login.Password == hash(input.Password) {
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
				"success": true,
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
				"success": false,
				"message": "Did you forget your own password?",
			})
			return
		}
	})
	r.DELETE("/", auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		userpub := UserPublic{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Where("id=?", id).Take(&userpub); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Delete(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Delete(&userpub); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"deleted": userpub,
		})
	})
	r.PATCH("/", auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var input User
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		userpub := UserPublic{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Where("id=?", id).Take(&userpub); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		user = User{
			ID:       user.ID,
			Name:     input.Name,
			Email:    input.Email,
			Password: hash(input.Password),
			Username: input.Username,
			NIK:      input.NIK,
			NIM:      input.NIM,
		}
		userpub = UserPublic{
			ID:       userpub.ID,
			Name:     input.Name,
			Email:    input.Email,
			Username: input.Username,
			NIK:      input.NIK,
			NIM:      input.NIM,
		}
		err := db.Model(&user).Updates(user)
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   err.Error.Error(),
			})
			return
		}
		ers := db.Model(&userpub).Updates(userpub)
		if ers.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   ers.Error.Error(),
			})
			return
		}
		if err := db.Where("id = ?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   err.Error.Error(),
			})
			return
		}
		if ers := db.Where("id = ?", id).Take(&userpub); ers.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   ers.Error.Error(),
			})
			return
		}
		if err.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		if ers.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    userpub,
			"input":   input,
		})
	})
	r.PATCH("/vaccine", auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var input UserVacc
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		vaksin := UserVacc{}
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
		if input.Dosis2 && !input.Booster {
			input.Dosis1 = true
			input.Booster = false
		}
		if input.Booster {
			input.Dosis1 = true
			input.Dosis2 = true
		}
		vaksin = UserVacc{
			ID:      vaksin.ID,
			Dosis1:  input.Dosis1,
			Dosis2:  input.Dosis2,
			Booster: input.Booster,
		}
		// err := db.Model(&vaksin).Updates(vaksin)
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
				"message": "User not found.",
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

func hash(input string) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	pass := hex.EncodeToString(hash.Sum(nil))
	return pass
}
