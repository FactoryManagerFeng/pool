package pool

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type Worker struct {
	job             chan *Job
	decNotify       chan bool
	stopCtx         context.Context
	stopCancelFunc  context.CancelFunc
	done            func()
	sleepSeconds    int64
	sleepCtx        context.Context
	sleepCancelFunc context.CancelFunc
	sleepNotify     chan bool
}

func (work *Worker) pushJob(job *Job) {
	work.job <- job
}

func (work *Worker) pushJobFunc(f JobFunc, args ...interface{}) {
	work.job <- &Job{
		f:    f,
		args: args,
	}
}

// 创建workers
func (work *Worker) createWorker(fu func()) {
	work.done = fu
	go work.doWork()
}

// 删除一个worker
func (work *Worker) deleteWorker() {
	work.decNotify <- true
}

// 停止整个workers
func (work *Worker) stopWorker() {
	work.stopCancelFunc()
}

// 休眠控制，当休眠到指定时间后，将时间重新设置为0
func (work *Worker) sleepControl() {
	defer work.done()
	for {
		select {
		case <-work.stopCtx.Done():
			return
		case isSleep := <-work.sleepNotify:
			if isSleep {
				work.sleepCancelFunc()
				time.Sleep(time.Second * time.Duration(work.sleepSeconds))
				work.sleepSeconds = 0
			}
		}
	}
}

// 休眠整个workers
func (work *Worker) sleepWorker(seconds int64) bool {
	// 判断是否设置成功
	if atomic.CompareAndSwapInt64(&work.sleepSeconds, 0, seconds) {
		work.sleepNotify <- true
		return true
	}
	return false
}

// worker具体处理逻辑
func (work *Worker) doWork() {
	defer work.done()
	for {
		select {
		case <-work.stopCtx.Done():
			return
		case <-work.sleepCtx.Done():
			if work.sleepSeconds > 0 {
				fmt.Println("job sleep,time:", work.sleepSeconds)
				time.Sleep(time.Second * time.Duration(work.sleepSeconds))
			}
		case flag := <-work.decNotify:
			if flag {
				return
			}
		case job := <-work.job:
			result := job.execute()
			fmt.Println("job done,result:", result)
		case <-time.After(15 * time.Second):
			fmt.Println("long time no message,exit")
			return
		}
	}
}
