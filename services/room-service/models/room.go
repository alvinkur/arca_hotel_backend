package models

type Room struct {
	ID           uint     `gorm:"primaryKey;column:id_room"    json:"id_room"`
	RoomNumber   string   `gorm:"not null;uniqueIndex"         json:"room_number"`
	RoomTypeID   uint     `gorm:"column:id_room_type;not null" json:"id_room_type"`
	Availability bool     `gorm:"default:true"                 json:"availability"`
	RoomType     RoomType `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
}

func (Room) TableName() string {
	return "room"
}
