package main

import (
	"fmt"

	"cloud-disk/app/common"
	"cloud-disk/internal/auth"
	"cloud-disk/internal/config"
	"cloud-disk/internal/log"
	"cloud-disk/internal/mysql"
	"cloud-disk/internal/server"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	serverAddr := config.AppCfg.ServerCfg.Host + ":" + config.AppCfg.ServerCfg.Port
	cloudDiskServer := server.NewServer(serverAddr)
	err = cloudDiskServer.Start()
	if err != nil {
		log.Error("start server error:%s", err)
		return
	}

	defer func() {
		err = cloudDiskServer.Close()
		if err != nil {
			log.Error("close cloud disk error:%s", err)
		}
		err = mysql.Close()
		if err != nil {
			log.Error("close mysql error:%s", err)
		}
		common.NewScheduledTask().Stop()
		log.Close()
		common.Close()
	}()
}

func initialize() error {
	err := config.InitConfig()
	if err != nil {
		return err
	}

	log.InitLog(&config.AppCfg.LogCfg)
	auth.InitAuth()

	//err = mysql.InitMySQL()
	//if err != nil {
	//	logs.Error("initialize MySQL error:%s", err)
	//	return err
	//}

	err = common.NewScheduledTask().StartScheduledTask()
	if err != nil {
		log.Error("start scheduled task error:%s", err)
		return err
	}

	err = common.InitGoroutinePool(config.AppCfg.ServerCfg.GoroutineNum)
	if err != nil {
		log.Error("initialize goroutine pool error:%s", err)
		return err
	}

	err = common.InitHttpClient()
	if err != nil {
		log.Error("initialize http client error:%s", err)
		return err
	}

	log.Info("success to initialize the cloud disk")
	return nil
}
