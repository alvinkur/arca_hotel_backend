package models

import "time"

type Chat struct {
	ID         uint      `gorm:"primaryKey;column:id_chat"   json:"id_chat"`
	CustomerID uint      `gorm:"column:id_customer;not null" json:"id_customer"`
	StaffID    uint      `gorm:"column:id_staff;not null"    json:"id_staff"`
	Date       time.Time `gorm:"autoCreateTime"              json:"date"`
}

func (Chat) TableName() string {
	return "chat"
}
