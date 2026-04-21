package models

type Customer struct {
	ID          uint   `gorm:"primaryKey;column:id_customer" json:"id_customer"`
	Name        string `gorm:"not null"                       json:"name"`
	Email       string `gorm:"uniqueIndex;not null"           json:"email"`
	Password    string `gorm:"not null"                       json:"-"`
	PhoneNumber string `json:"phone_number"`
}

func (Customer) TableName() string {
	return "customer"
}
