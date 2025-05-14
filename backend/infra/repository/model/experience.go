package model

type Experience struct {
	ID    int    `gorm:"primaryKey"`
	Title string `gorm:"not null unique"`
}

func (e *Experience) TableName() string {
	return "experiences"
}
