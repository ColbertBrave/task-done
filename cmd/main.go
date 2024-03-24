package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/task-done/infrastructure/auth"
	"github.com/task-done/infrastructure/client"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/influxdb"
	"github.com/task-done/infrastructure/log"
	"github.com/task-done/infrastructure/mysql"
	"github.com/task-done/infrastructure/pool"
	"github.com/task-done/infrastructure/server"
	"github.com/task-done/infrastructure/task"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run()
	log.System("successfully start the server")

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

	log.System("successfully initialize the service")
	return nil
}

func finalize() {
	log.Info("exit the process!")
	server.Close()
	if err := mysql.Close(); err != nil {
		log.System("close mysql error|%s", err)
	}
	task.NewScheduledTask().Stop()
	log.Close()
}
