package models

import "time"

type Payment struct {
	ID           uint      `gorm:"primaryKey;column:id_payment"    json:"id_payment"`
	BookingID    uint      `gorm:"column:id_booking;not null"      json:"id_booking"`
	TotalPayment float64   `gorm:"not null"                        json:"total_payment"`
	Method       string    `gorm:"not null"                        json:"method"`
	Date         time.Time `json:"date"`
	Status       string    `gorm:"default:pending"                 json:"status"`

	Booking Booking `gorm:"foreignKey:BookingID;references:ID" json:"booking,omitempty"`
}

func (Payment) TableName() string {
	return "payment"
}
