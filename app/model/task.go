package model

type Task struct {
	Model
	Name        string `gorm:"not_null" json:"task"`
	Description string `json:"Description"`
}

func (t *Task) TableName() string {
	return "task_info"
}
