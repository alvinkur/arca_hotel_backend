package models

type Staff struct {
	ID       uint   `gorm:"primaryKey;column:id_staff" json:"id_staff"`
	Name     string `gorm:"not null"                   json:"name"`
	Email    string `gorm:"uniqueIndex;not null"       json:"email"`
	Password string `gorm:"not null"                   json:"-"`
}

func (Staff) TableName() string {
	return "staff"
}
