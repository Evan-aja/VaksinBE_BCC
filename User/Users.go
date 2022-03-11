package User

import (
	"VaksinBE_BCC/Auth"
	"VaksinBE_BCC/Vaccine"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
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
	// show logged in user profile
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
				"id":            user.ID,
				"name":          user.Name,
				"email":         user.Email,
				"handphone":     user.Handphone,
				"nim":           user.NIM,
				"nik":           user.NIK,
				"gender":        user.Gender,
				"tanggal_lahir": user.TanggalLahir,
				"vaksinasi":     vacc,
			},
		})
	})
	// search for other user profile by name or email
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
	// normal registration without google. needs Name, Email, Handphone, and Password

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
		}
		registPublic := UserPublic{}
		registVacc := Vaccine.Vaccine{
			Dosis1:  false,
			Dosis2:  false,
			Booster: false,
		}
		regist.TanggalLahir, _ = time.Parse("2006-01-02T15:04:05Z", "0001-01-01T00:00:00Z")
		if err := db.Create(&regist); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with user creation",
				"error":   err.Error.Error(),
			})
			return
		}
		registPublic = UserPublic{
			PubID: regist.ID,
			Name:  input.Name,
			Email: input.Email,
		}
		if err := db.Create(&registPublic); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with public user creation",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Create(&registVacc); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went  on vaccine setter",
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
				"message": "Something went wrong on proof vaccine setter",
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
	// normal login without google. needs Email/Handphone, and Password
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
					"message": "Phone does not exist",
					"error":   err.Error.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Email does not exist",
				"error":   err.Error.Error(),
			})
			return
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

	// google registration, needs a button on web page for redirect and subsequently, logged in
	r.GET("/google", Auth.GInit)
	r.GET("/google/callback", func(c *gin.Context) {
		var a = Auth.GCallback(c)
		var b = []byte(a)
		var goog Goog
		var user User
		if err := json.Unmarshal(b, &goog); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"err":     err.Error(),
			})
			return
		}
		user = User{}
		if err := db.Where("email=?", goog.Email).Take(&user); err.Error != nil {
			user = User{
				Name:     goog.Name,
				Email:    goog.Email,
				Password: hash(goog.Sub),
			}
			registPublic := UserPublic{}
			registVacc := Vaccine.Vaccine{
				Dosis1:  false,
				Dosis2:  false,
				Booster: false,
			}
			user.TanggalLahir, _ = time.Parse("2006-01-02 15:04", "0001-01-01T00:00:00Z")
			if err := db.Create(&user); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Something went wrong with user creation",
					"error":   err.Error.Error(),
				})
				return
			}
			registPublic = UserPublic{
				PubID: user.ID,
				Name:  goog.Name,
				Email: goog.Email,
			}
			if err := db.Create(&registPublic); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Something went wrong with public user creation",
					"error":   err.Error.Error(),
				})
				return
			}
			if err := db.Create(&registVacc); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Something went wrong with vaccine setter",
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
					"message": "Something went wrong with proof vaccine setter",
					"error":   err.Error.Error(),
				})
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "no account was found",
				"status":  "creating new account",
			})
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"id":  user.ID,
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
				"email":     user.Email,
				"handphone": user.Handphone,
				"token":     strToken,
			},
		})
	})

	// deletes user data, universal, only needs to be logged in. FE should gives prompt to ask wether the account deletion is intentional or not.
	r.DELETE("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		vacc := Vaccine.Vaccine{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
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
		if err := db.Where("ID=?", vacc.ID).Delete(&vacc); err.Error != nil {
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
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"deleted": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	})

	// updates user data (universal). can be filled all willy nilly. users can change Name, Email, Handphone, and Password. FE should restrict Email, however, since if an invalid email is written by the user, it'll make their account inaccessible
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
			ID:           user.ID,
			Name:         input.Name,
			Email:        input.Email,
			Password:     hash(input.Password),
			Handphone:    input.Handphone,
			TanggalLahir: input.TanggalLahir,
			NIK:          input.NIK,
			NIM:          input.NIM,
			Gender:       input.Gender,
		}
		userpub = UserPublic{
			PubID:  user.ID,
			Name:   input.Name,
			Email:  input.Email,
			Gender: input.Gender,
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
				"id":    userpub.PubID,
				"name":  userpub.Name,
				"email": userpub.Email,
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
