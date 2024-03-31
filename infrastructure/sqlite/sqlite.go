package sqlite

import (
	"fmt"

	"github.com/task-done/infrastructure/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
) 

func Init() error {
	path := config.GetConfig().SQLite.Path
	sqlPointer, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("init sqlite error|%s|%s",err, path)
	}

	db = sqlPointer
	return nil
}

func GetDB() *gorm.DB {
	if db == nil {
		Init()
	}
	return db
}

func AutoMigrate(objects ...interface{}) error {
	return db.AutoMigrate(objects...)
} 