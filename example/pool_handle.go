package logic

import (
	"code.piupiu.com/pool"
	"fmt"
	"time"
)

var name = "test"

func PoolStart() {
	newPool, err := pool.NewPool(name, 100, 1)
	if err != nil {
		return
	}
	newPool.Start()

	for i := 0; i < 10000; i++ {
		fmt.Println("i", i)
		newPool.PushJobFunc(func(args ...interface{}) pool.State {
			arg := args[0].([]interface{})
			fmt.Println("args", args, "arg", arg)
			time.Sleep(time.Second)
			if arg[0].(int) > 100 {
				return pool.StateOk
			}
			return pool.StateErr
		}, i)
	}
	newPool.Stop()
	return
}

func PoolStop() {
	p := pool.Get(name)
	if p == nil {
		fmt.Println("pool is empty", p)
		return
	}
	p.Stop()
	fmt.Println("pool is stop", p)
	return
}

func PoolAdd() {
	p := pool.Get(name)
	if p == nil {
		fmt.Println("pool is empty", p)
		return
	}
	p.IncWorker(10)
	return
}

func PoolDes() {
	p := pool.Get(name)
	if p == nil {
		fmt.Println("pool is empty", p)
		return
	}
	p.DecWorker(19)
	return
}

func PoolSleep() {
	p := pool.Get(name)
	if p == nil {
		fmt.Println("pool is empty", p)
		return
	}
	p.Sleep(5)
}
