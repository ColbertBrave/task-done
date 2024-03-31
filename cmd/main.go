package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/task-done/app/api"
	"github.com/task-done/app/model"
	"github.com/task-done/infrastructure/auth"
	"github.com/task-done/infrastructure/client"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/influxdb"
	"github.com/task-done/infrastructure/log"
	"github.com/task-done/infrastructure/pool"
	"github.com/task-done/infrastructure/server"
	"github.com/task-done/infrastructure/sqlite"
	"github.com/task-done/infrastructure/task"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 自动建表
	if err := sqlite.AutoMigrate(&model.Task{}, &model.User{}); err != nil {
		log.System("auto miragte err|%s",err)
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
	err := config.Init()
	if err != nil {
		return err
	}

	log.Init()
	auth.Init()
	influxdb.Init()

	err = sqlite.Init()
	if err != nil {
		log.Error("initialize sqlite err|%s", err)
		return err
	}

	err = task.NewOnce().StartScheduledTask()
	if err != nil {
		log.Error("start scheduled task err|%s", err)
		return err
	}

	err = pool.Init()
	if err != nil {
		log.Error("initialize goroutine pool err|%s", err)
		return err
	}

	client.Init()
	server.Init()
	api.Init()

	log.System("successfully initialize the service")
	return nil
}

func finalize() {
	log.Info("exit the process!")
	server.Close()
	task.NewOnce().Stop()
	log.Close()
}
