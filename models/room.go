package models

type Room struct {
	ID           uint    `gorm:"primaryKey;column:id_room" json:"id_room"`
	Type         string  `gorm:"not null"                  json:"type"`
	Price        float64 `gorm:"not null"                  json:"price"`
	Availability bool    `gorm:"default:true"              json:"availability"`
}

func (Room) TableName() string {
	return "room"
}
