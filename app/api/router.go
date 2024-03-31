package api

import (
	"github.com/task-done/app/service/task"
	"github.com/task-done/infrastructure/constants"
	"github.com/task-done/infrastructure/server"
)

func Init()  {
	server.GET(constants.PrefixURL+"/task/info", task.GetTaskInfo)
}
