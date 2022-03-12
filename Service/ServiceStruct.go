package Service

import (
	"VaksinBE_BCC/User"
	"time"
)

type Swab struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Type     string     `gorm:"uniqueIndex;notnull;size:256" json:"type"`
	Cost     uint       `json:"cost"`
	AdCost   uint       `json:"admin_cost"`
	Schedule []Schedule `gorm:"many2many:swab_schedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AddSwab struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Type       string `gorm:"uniqueIndex;notnull;size:256" json:"type"`
	Cost       uint   `json:"cost"`
	AdCost     uint   `json:"admin_cost"`
	ScheduleID []uint `json:"id_schedule"`
}

type Poly struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Type        string     `gorm:"notnull;size:256" json:"type"`
	Description string     `gorm:"notnull" json:"description"`
	Cost        uint       `json:"cost"`
	AdCost      uint       `json:"admin_cost"`
	Schedule    []Schedule `gorm:"many2many:poly_schedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AddPoly struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Type        string `gorm:"notnull;size:256" json:"type"`
	Description string `gorm:"notnull" json:"description"`
	Cost        uint   `json:"cost"`
	AdCost      uint   `json:"admin_cost"`
	ScheduleID  []uint `json:"id_schedule"`
}

type Schedule struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Hour string `json:"hour"`
}

type TransactionSwab struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User.User `gorm:"foreignKey:UserID"`
	Name         string    `gorm:"notNull" json:"name"`
	Email        string    `gorm:"notNull;size:256" json:"email"`
	Handphone    string    `json:"handphone"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	NIK          string    `json:"nik"`
	NIM          string    `json:"nim"`
	Gender       string    `json:"gender"`
	Type         string    `gorm:"size:256" json:"type"`
	Date         time.Time `gorm:"notNull" json:"date"`
	Cost         uint      `json:"cost"`
	AdCost       uint      `json:"admin_cost"`
	IDSchedule   uint      `json:"id_schedule"`
	Schedule     Schedule  `gorm:"foreignKey:IDSchedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Paid         bool      `json:"status_pay"`
	Receipt      string    `json:"receipt"`
}

type TransSwabUp struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Paid    bool   `json:"status_pay"`
	Receipt string `json:"receipt"`
}

type TransactionPoly struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User.User `gorm:"foreignKey:UserID"`
	Name         string    `gorm:"notNull" json:"name"`
	Email        string    `gorm:"notNull;size:256" json:"email"`
	Handphone    string    `json:"handphone"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	NIK          string    `json:"nik"`
	NIM          string    `json:"nim"`
	Gender       string    `json:"gender"`
	Type         string    `gorm:"size:256" json:"type"`
	Description  string    `gorm:"notnull" json:"description"`
	Date         time.Time `gorm:"notNull" json:"date"`
	Cost         uint      `json:"cost"`
	AdCost       uint      `json:"admin_cost"`
	IDSchedule   uint      `json:"id_schedule"`
	Schedule     Schedule  `gorm:"foreignKey:IDSchedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Paid         bool      `json:"status_pay"`
	Receipt      string    `json:"receipt"`
}

type TransPolyUp struct {
	ID      uint   `json:"id"`
	Paid    bool   `json:"status_pay"`
	Receipt string `json:"receipt"`
}
