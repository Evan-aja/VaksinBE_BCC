package User

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	NIK      string `gorm:"uniqueIndex;notNull" json:"nik"`
	NIM      string `gorm:"uniqueIndex;notNull" json:"nim"`
}
type UserPublic struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	NIK      string `gorm:"uniqueIndex;notNull" json:"nik"`
	NIM      string `gorm:"uniqueIndex;notNull" json:"nim"`
}

type UserVacc struct {
	ID      uint `gorm:"primarykey" json:"id"`
	Dosis1  bool `gorm:"notNull" json:"dosis1"`
	Dosis2  bool `gorm:"notNull" json:"dosis2"`
	Booster bool `gorm:"notNull" json:"booster"`
}

type UserRegister struct {
	Name     string `gorm:"notNull" json:"name"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	NIK      string `gorm:"uniqueIndex;notNull" json:"nik"`
	NIM      string `gorm:"uniqueIndex;notNull" json:"nim"`
}

type UserLogin struct {
	Username string `gorm:"uniqueIndex;notNull" json:"username"`
	Email    string `gorm:"uniqueIndex;notNull" json:"email"`
	Password string `gorm:"notNull" json:"password"`
}
