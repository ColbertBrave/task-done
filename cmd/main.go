package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloud-disk/infrastructure/auth"
	"github.com/cloud-disk/infrastructure/client"
	"github.com/cloud-disk/infrastructure/config"
	"github.com/cloud-disk/infrastructure/influxdb"
	"github.com/cloud-disk/infrastructure/log"
	"github.com/cloud-disk/infrastructure/mysql"
	"github.com/cloud-disk/infrastructure/pool"
	"github.com/cloud-disk/infrastructure/server"
	"github.com/cloud-disk/infrastructure/task"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	for range ctx.Done() {
		finalize()
		stop()
	}
}

func initialize() error {
	err := config.InitConfig()
	if err != nil {
		return err
	}

	log.InitLog()
	auth.InitAuth()
	influxdb.InitInfluxdb()

	//err = mysql.InitMySQL()
	//if err != nil {
	//	logs.Error("initialize MySQL error: %s", err)
	//	return err
	//}

	// err = mongodb.InitMongoDB()
	// if err != nil {
	// 	log.Error("initialize MongoDB error: %s", err)
	// 	return err
	// }

	err = task.NewScheduledTask().StartScheduledTask()
	if err != nil {
		log.Error("start scheduled task error: %s", err)
		return err
	}

	err = pool.InitGoroutinePool()
	if err != nil {
		log.Error("initialize goroutine pool error: %s", err)
		return err
	}

	client.InitHttpClient()
	server.InitServer()

	log.Info("success to initialize the cloud disk")
	return nil
}

func finalize() {
	log.Info("exit the process!")
	server.Close()
	if err := mysql.Close(); err != nil {
		log.Error("close mysql error|%s", err)
	}
	task.NewScheduledTask().Stop()
	log.Close()
}
