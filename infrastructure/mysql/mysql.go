package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/cloud-disk/infrastructure/config"
	"github.com/cloud-disk/infrastructure/log"
)

var MySqlTables *Tables

type Tables struct {
	dbPointer *gorm.DB
}

func InitMySQL() error {
	userName := config.GetConfig().MySQLCfg.UserName
	password := config.GetConfig().MySQLCfg.Password
	host := config.GetConfig().MySQLCfg.Host
	port := config.GetConfig().MySQLCfg.Port
	database := config.GetConfig().MySQLCfg.Database

	dbURL := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		userName, password, host, port, database)
	dbConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	dbPointer, err := gorm.Open(mysql.Open(dbURL), &dbConfig)
	if err != nil {
		log.Error("init mysql error:%s", err)
		return err
	}

	MySqlTables = &Tables{}

	err = dbPointer.AutoMigrate()
	if err != nil {
		log.Error("auto migrate tables error:%s", err)
		return err
	}

	db, err := dbPointer.DB()
	if err != nil {
		log.Error("set max connection time error:%s", err)
		return err
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	return nil
}

func Close() error {
	db, err := MySqlTables.dbPointer.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
