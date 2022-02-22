package User

import "github.com/gin-gonic/gin"

var r *gin.Engine

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	NIK      string `gorm:"uniqueIndex;notNull"`
	NIM      string `gorm:"uniqueIndex;notNull"`
}

type UserRegister struct {
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	NIK      string `gorm:"uniqueIndex;notNull"`
	NIM      string `gorm:"uniqueIndex;notNull"`
}

type UserLogin struct {
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
}
