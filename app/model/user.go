package model

type User struct {
	Model
	Name        string `gorm:"not_null" json:"task"`
	Mail        string `json:"mail"`
	Description string `json:"Description"`
}

func (u *User) TableName() string {
	return "user_info"
}
