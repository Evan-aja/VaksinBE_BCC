package User

type User struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"notNull" json:"name"`
	Email     string `gorm:"uniqueIndex;notNull" json:"email"`
	Password  string `gorm:"notNull" json:"password"`
	Handphone string `gorm:"uniqueIndex;notNull" json:"handphone"`
}
type UserPublic struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"notNull" json:"name"`
	Email     string `gorm:"uniqueIndex;notNull" json:"email"`
	Handphone string `gorm:"uniqueIndex;notNull" json:"handphone"`
}

type UserRegister struct {
	Name      string `gorm:"notNull" json:"name"`
	Email     string `gorm:"uniqueIndex;notNull" json:"email"`
	Password  string `gorm:"notNull" json:"password"`
	Handphone string `gorm:"uniqueIndex;notNull" json:"handphone"`
}

type UserLogin struct {
	Handphone string `gorm:"uniqueIndex;notNull" json:"handphone"`
	Email     string `gorm:"uniqueIndex;notNull" json:"email"`
	Password  string `gorm:"notNull" json:"password"`
}
