package models

type RevenueReport struct {
	ID           uint    `gorm:"primaryKey;column:id_revenue" json:"id_revenue"`
	Period       string  `gorm:"not null"                     json:"period"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalBooking int     `json:"total_booking"`
	TotalReview  int     `json:"total_review"`
	DetailIncome string  `json:"detail_income"`
}

func (RevenueReport) TableName() string {
	return "revenue_report"
}
