package models

import "time"

type Booking struct {
	ID            uint      `gorm:"primaryKey;column:id_booking" json:"id_booking"`
	CustomerID    uint      `gorm:"column:id_customer;not null"  json:"id_customer"`
	RoomID        uint      `gorm:"column:id_room;not null"      json:"id_room"`
	DateIn        time.Time `gorm:"not null"                     json:"date_in"`
	DateOut       time.Time `gorm:"not null"                     json:"date_out"`
	TotalPayment  float64   `json:"total_payment"`
	StatusPayment string    `gorm:"default:pending"              json:"status_payment"`
}

func (Booking) TableName() string {
	return "booking"
}
