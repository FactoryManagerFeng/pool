package pool

import (
	"context"
	"fmt"
)

type Worker struct {
	job            chan *Job
	dec            chan bool
	stopCtx        context.Context
	stopCancelFunc context.CancelFunc
	done           func()
}

func NewWorker() *Worker {
	return &Worker{
		job: make(chan *Job),
		dec: make(chan bool),
	}
}

func (work *Worker) PushJob(job *Job) {
	work.job <- job
}

func (work *Worker) PushJobFunc(f JobFunc, args ...interface{}) {
	work.job <- &Job{
		f:    f,
		args: args,
	}
}

// 创建一个worker
func (work *Worker) createWorker(fu func()) {
	work.done = fu
	work.stopCtx, work.stopCancelFunc = context.WithCancel(context.Background())

	go work.doWork()
}

// 删除一个worker
func (work *Worker) deleteWorker() {
	work.dec <- true
}

// 停止整个workers
func (work *Worker) stopAllWorker() {
	work.stopCancelFunc()
}

// worker具体处理逻辑
func (work *Worker) doWork() {
	defer work.done()
	for {
		select {
		case <-work.stopCtx.Done():
			return
		case flag := <-work.dec:
			if flag {
				return
			}
		case job := <-work.job:
			state := job.execute()
			fmt.Println("job done,state:", state)
		}
	}
}
