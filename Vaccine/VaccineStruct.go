package Vaccine

type Vaccine struct {
	ID      uint `gorm:"primarykey" json:"id"`
	Dosis1  bool `gorm:"notNull" json:"dosis1"`
	Dosis2  bool `gorm:"notNull" json:"dosis2"`
	Booster bool `gorm:"notNull" json:"booster"`
}

type VaccProof struct {
	IDVaccine uint
	Vaccine   Vaccine `gorm:"ForeignKey:IDVaccine;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"idbukti"`
	Dosis1    string  `gorm:"null" json:"buktidosis1"`
	Dosis2    string  `gorm:"null" json:"buktidosis2"`
	Booster   string  `gorm:"null" json:"buktibooster"`
}
type VaccineIn struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Dosis1  bool   `gorm:"notNull" json:"dosis1"`
	Dosis2  bool   `gorm:"notNull" json:"dosis2"`
	Booster bool   `gorm:"notNull" json:"booster"`
	Bukti1  string `gorm:"null" json:"buktidosis1"`
	Bukti2  string `gorm:"null" json:"buktidosis2"`
	Bukti3  string `gorm:"null" json:"buktibooster"`
}
