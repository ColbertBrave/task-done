package pool

import (
	"github.com/task-done/app/types/result"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/log"

	"github.com/panjf2000/ants/v2"
)

var Pool *GoroutinePool

type GoroutinePool struct {
	goroutinePool *ants.Pool
}

func Init() error {
	num := config.GetConfig().Server.PoolGoroutineNum
	if num <= 0 {
		return result.ErrInvalidInputParam
	}

	pool, err := ants.NewPool(num)
	if err != nil {
		log.Error("fail to new a goroutine pool, err:%s", err)
		return err
	}
	Pool = &GoroutinePool{pool}
	return nil
}

func (g *GoroutinePool) Submit(task func()) {
	g.goroutinePool.Submit(task)
}

func Close() {
	if Pool != nil {
		Pool.goroutinePool.Release()
	}
}
