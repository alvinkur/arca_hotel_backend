package models

type RoomType struct {
	ID          uint    `gorm:"primaryKey;column:id_room_type" json:"id_room_type"`
	Name        string  `gorm:"not null;uniqueIndex"           json:"name"`
	Price       float64 `gorm:"not null"                       json:"price"`
	Description string  `json:"description"`
}

func (RoomType) TableName() string {
	return "room_type"
}
