package Vaccine

type Vaccine struct {
	ID      uint `gorm:"primarykey" json:"id"`
	Dosis1  bool `gorm:"notNull" json:"dosis1"`
	Dosis2  bool `gorm:"notNull" json:"dosis2"`
	Booster bool `gorm:"notNull" json:"booster"`
}
