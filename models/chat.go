package models

import "time"

type Chat struct {
	ID         uint      `gorm:"primaryKey;column:id_chat"           json:"id_chat"`
	CustomerID uint      `gorm:"column:id_customer;not null"         json:"id_customer"`
	StaffID    uint      `gorm:"column:id_staff;not null"            json:"id_staff"`
	Chat       string    `gorm:"not null"                            json:"chat"`
	Date       time.Time `gorm:"autoCreateTime"                      json:"date"`

	Customer Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Staff    Staff    `gorm:"foreignKey:StaffID;references:ID"    json:"staff,omitempty"`
}

func (Chat) TableName() string {
	return "chat"
}
