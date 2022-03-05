package User

import (
	"time"
)

// regular user data
type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `gorm:"notNull" json:"name"`
	Email        string    `gorm:"uniqueIndex;notNull" json:"email"`
	Password     string    `gorm:"notNull" json:"password"`
	Handphone    string    `json:"handphone"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	NIK          string    `json:"nik"`
	NIM          string    `json:"nim"`
	Gender       string    `json:"gender"`
}

// publicly available user data
type UserPublic struct {
	PubID  uint   `json:"pubid"`
	User   User   `gorm:"ForeignKey:PubID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name   string `gorm:"notNull" json:"name"`
	Email  string `gorm:"uniqueIndex;notNull" json:"email"`
	Gender string `json:"gender"`
}

// google data for login and stuff
type Goog struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	EmailVerif bool   `json:"email_verified"`
	Locale     string `json:"locale"`
}

// regular register
type UserRegister struct {
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
}

// regular login
type UserLogin struct {
	Handphone string `gorm:"uniqueIndex;notNull" json:"handphone"`
	Email     string `gorm:"uniqueIndex;notNull" json:"email"`
	Password  string `gorm:"notNull" json:"password"`
}
