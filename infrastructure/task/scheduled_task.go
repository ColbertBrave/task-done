package task

import (
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/log"
)

var task *ScheduledTask

type ScheduledTask struct {
	c *cron.Cron
}

var once sync.Once

func NewOnce() *ScheduledTask {
	once.Do(func() {
		task = &ScheduledTask{
			c: cron.New(),
		}
	})
	return task
}

// StartScheduledTask 启动定时任务
func (s *ScheduledTask) StartScheduledTask() error {
	_, err := s.c.AddFunc(config.GetConfig().TaskTime.QueryTime, test)
	if err != nil {
		log.Error("add scheduled task error:%s", err)
		return err
	}

	s.c.Start()
	return nil
}

func (s *ScheduledTask) Stop() {
	s.c.Stop()
}

func test() {

}
