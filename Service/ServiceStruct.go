package Service

import (
	"VaksinBE_BCC/User"
	"time"
)

type Swab struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Type     string     `gorm:"uniqueIndex;size:256" json:"type"`
	Cost     uint       `json:"cost"`
	AdCost   uint       `json:"admin_cost"`
	Schedule []Schedule `gorm:"many2many:swab_schedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AddSwab struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Type       string `gorm:"uniqueIndex;size:256" json:"type"`
	Cost       uint   `json:"cost"`
	AdCost     uint   `json:"admin_cost"`
	ScheduleID []uint `json:"id_schedule"`
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
	// Schedule     []Schedule `json:"schedule"`
	// IDSwab       uint       `json:"id_swab"`
	// Swab         Swab       `gorm:"foreignKey:IDSwab"`
	Type       string    `gorm:"uniqueIndex;size:256" json:"type"`
	Date       time.Time `gorm:"notNull" json:"date"`
	Cost       uint      `json:"cost"`
	AdCost     uint      `json:"admin_cost"`
	IDSchedule uint      `json:"id_schedule"`
	Schedule   Schedule  `gorm:"foreignKey:IDSchedule;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Paid       bool      `json:"status_pay"`
}
