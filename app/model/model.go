package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/task-done/infrastructure/sqlite"
	"gorm.io/gorm"
)

type Model struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) Create() error {
	return sqlite.GetDB().Create(m).Error
}

func (m *Model) First() error {
	return sqlite.GetDB().First(m).Error
}

func (m *Model) Last() error {
	return sqlite.GetDB().Last(m).Error
}

func (m *Model) Find() error {
	return sqlite.GetDB().Find(m).Error
}

func (m *Model) Update() error {
	return sqlite.GetDB().Save(m).Error
}

func (m *Model) Delete() error {
	return sqlite.GetDB().Delete(m).Error
}

func (m *Model) Raw(sql string) error {
	return sqlite.GetDB().Raw(sql).Scan(m).Error
}

func (m *Model) BeforeCreate(tx *gorm.DB) {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now()
}

func (m *Model) BeforeUpdate(tx *gorm.DB) {
	m.UpdatedAt = time.Now()
}