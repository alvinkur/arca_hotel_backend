package models

type Owner struct {
	ID       uint   `gorm:"primaryKey;column:id_owner" json:"id_owner"`
	Name     string `gorm:"not null"                   json:"name"`
	Email    string `gorm:"uniqueIndex;not null"       json:"email"`
	Password string `gorm:"not null"                   json:"-"`
}

func (Owner) TableName() string {
	return "owner"
}
