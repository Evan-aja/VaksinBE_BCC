package Vaccine

// vaccine boolean (true or false)
type Vaccine struct {
	ID      uint `gorm:"primarykey" json:"id"`
	Dosis1  bool `gorm:"notNull" json:"dosis1"`
	Dosis2  bool `gorm:"notNull" json:"dosis2"`
	Booster bool `gorm:"notNull" json:"booster"`
}

// proof of vaccine (links to images uploaded)
type VaccProof struct {
	IDVaccine uint    `json:"idbukti"`
	Vaccine   Vaccine `gorm:"ForeignKey:IDVaccine;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Dosis1    string  `gorm:"null" json:"buktidosis1"`
	Dosis2    string  `gorm:"null" json:"buktidosis2"`
	Booster   string  `gorm:"null" json:"buktibooster"`
}

// aggregate for previous 2 structures
type VaccineIn struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Dosis1  bool   `gorm:"notNull" json:"dosis1"`
	Dosis2  bool   `gorm:"notNull" json:"dosis2"`
	Booster bool   `gorm:"notNull" json:"booster"`
	Bukti1  string `gorm:"null" json:"buktidosis1"`
	Bukti2  string `gorm:"null" json:"buktidosis2"`
	Bukti3  string `gorm:"null" json:"buktibooster"`
}
