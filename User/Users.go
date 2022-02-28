package User

import (
	"VaksinBE_BCC/Auth"
	"VaksinBE_BCC/Vaccine"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/user")
	r.GET("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		vacc := Vaccine.VaccProof{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Where("id_vaccine=?", id).Preload("Vaccine").Take(&vacc); err.Error != nil {
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
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"handphone": user.Handphone,
				"vaksinasi": vacc,
			},
		})
	})
	r.GET("/search", func(c *gin.Context) {
		name, nameExists := c.GetQuery("name")
		email, emailExists := c.GetQuery("email")
		handphone, handphoneExists := c.GetQuery("handphone")
		if !nameExists && !handphoneExists && !emailExists {
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
		if handphoneExists {
			dbtmp = dbtmp.Where("handphone LIKE ?", "%"+handphone+"%")
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
					"name":      name,
					"email":     email,
					"handphone": handphone,
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
			Name:      input.Name,
			Email:     input.Email,
			Password:  hash(input.Password),
			Handphone: input.Handphone,
		}
		registPublic := UserPublic{}
		registVacc := Vaccine.Vaccine{
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
		registPublic = UserPublic{
			PubID:     regist.ID,
			Name:      input.Name,
			Email:     input.Email,
			Handphone: input.Handphone,
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
		ProofVacc := Vaccine.VaccProof{
			IDVaccine: registVacc.ID,
		}
		if err := db.Create(&ProofVacc); err.Error != nil {
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
			"data":    registPublic,
		})
	})
	r.GET("/google", Auth.GInit)
	r.GET("/google/callback", Auth.GCallback, func(c *gin.Context) {
		// c.BindJSON()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    "ID",
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
			if err = db.Where("handphone=?", input.Handphone).Take(&login); err.Error != nil {
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
			godotenv.Load(".env")
			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
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
					"email":     login.Email,
					"handphone": login.Handphone,
					"token":     strToken,
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
	r.DELETE("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		userpub := UserPublic{}
		vacc := Vaccine.Vaccine{}
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
		if err := db.Where("id=?", id).Take(&vacc); err.Error != nil {
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
		if err := db.Delete(&vacc); err.Error != nil {
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
	r.PATCH("/", Auth.Authorization(), func(c *gin.Context) {
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
		user = User{
			ID:        user.ID,
			Name:      input.Name,
			Email:     input.Email,
			Password:  hash(input.Password),
			Handphone: input.Handphone,
		}
		userpub = UserPublic{
			PubID:     user.ID,
			Name:      input.Name,
			Email:     input.Email,
			Handphone: input.Handphone,
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
		ers := db.Model(&userpub).Where("pub_id=?", user.ID).Updates(userpub)
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
		if ers := db.Where("pub_id = ?", id).Preload("User").Take(&userpub); ers.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   ers.Error.Error(),
			})
			return
		}
		if err.RowsAffected < 1 || ers.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No data has been changed.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data": gin.H{
				"id":        userpub.PubID,
				"name":      userpub.Name,
				"email":     userpub.Email,
				"handphone": userpub.Handphone,
			},
			"input": gin.H{
				"id":        input.ID,
				"name":      input.Name,
				"email":     input.Email,
				"handphone": input.Handphone,
			},
		})
	})
}

func hash(input string) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	pass := hex.EncodeToString(hash.Sum(nil))
	return pass
}
