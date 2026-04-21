package models

type Review struct {
	ID         uint   `gorm:"primaryKey;column:id_review"               json:"id_review"`
	CustomerID uint   `gorm:"column:id_customer;not null"               json:"id_customer"`
	RoomID     uint   `gorm:"column:id_room;not null"                   json:"id_room"`
	Rating     int    `gorm:"check:rating >= 1 AND rating <= 5;not null" json:"rating"`
	Comment    string `json:"comment"`

	Customer Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Room     Room     `gorm:"foreignKey:RoomID;references:ID"     json:"room,omitempty"`
}

func (Review) TableName() string {
	return "review"
}
