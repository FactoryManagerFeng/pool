package logic

import (
	"code.qschou.com/vip/go_vip_common/dlog"
	"code.qschou.com/vip/go_vip_common/pool"
	"fmt"
	"time"
)

var name = "test"

func PoolStart() {
	worker := pool.NewWorker()
	newPool := pool.NewPool(name, 100, 1, worker)
	newPool.Start()

	for i := 0; i < 10000; i++ {
		fmt.Println("i", i)
		worker.PushJobFunc(func(args ...interface{}) pool.State {
			fmt.Println("args", args)
			time.Sleep(1 * time.Second)
			return pool.StateOk
		}, i)
	}
	return
}

func PoolStop() {
	p := pool.Get(name)
	if p == nil {
		dlog.Info("pool is empty", p)
		return
	}
	p.Stop()
	dlog.Info("pool is stop", p)
	return
}

func PoolAdd() {
	p := pool.Get(name)
	if p == nil {
		dlog.Info("pool is empty", p)
		return
	}
	p.IncWorker(10)
	return
}

func PoolDes() {
	p := pool.Get(name)
	if p == nil {
		dlog.Info("pool is empty", p)
		return
	}
	p.DecWorker(19)
	return
}
