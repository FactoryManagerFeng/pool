package pool

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

const (
	Capacity = 500
	Running  = 1
)

type Pool struct {
	name     string         //线程池名称，不同线程池有不同的名字
	capacity int            //最大线程数
	running  int            //当前工作线程数
	worker   *Worker        //任务
	wg       sync.WaitGroup // wait
	lock     sync.Mutex     //锁
}

// 初始化协程
func NewPool(name string, capacity, running int) (*Pool, error) {
	if p, ok := poolMap[strings.ToLower(name)]; ok {
		return p, ErrPoolHaveSameName
	}
	if capacity > Capacity {
		capacity = Capacity
	}
	if running < Running {
		running = Running
	}
	pool := &Pool{
		name:     name,
		capacity: capacity,
		running:  running,
		worker: &Worker{
			job:         make(chan *Job),
			decNotify:   make(chan bool),
			sleepNotify: make(chan bool),
		},
	}
	return pool, nil
}

// 添加工作线程数 最多添加到设置的最大线程数
func (p *Pool) IncWorker(num int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	surplus := p.capacity - p.running
	if num > surplus {
		num = surplus
	}
	for i := 0; i < num; i++ {
		p.running++
		p.wg.Add(1)
		p.worker.createWorker(func() { p.wg.Done() })
	}
	fmt.Println(p)
	return
}

// 减少工作线程数 最多减少到剩余一个工作线程
func (p *Pool) DecWorker(num int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if num >= p.running {
		num = p.running - 1
	}
	for i := 0; i < num; i++ {
		p.running--
		p.worker.deleteWorker()
	}
	return
}

// 开始工作
func (p *Pool) Start() bool {
	defer p.lock.Unlock()
	p.lock.Lock()
	if p.running > p.capacity || p.running <= 0 {
		return false
	}

	p.worker.stopCtx, p.worker.stopCancelFunc = context.WithCancel(context.Background())
	p.worker.sleepCtx, p.worker.sleepCancelFunc = context.WithCancel(context.Background())
	p.wg.Add(p.running + 1)
	for i := 0; i < p.running; i++ {
		p.worker.createWorker(func() { p.wg.Done() })
	}
	register(p)
	go p.worker.sleepControl()
	return true
}

// 停止工作
func (p *Pool) Stop() {
	defer p.lock.Unlock()
	p.lock.Lock()

	p.worker.stopWorker()
	p.wg.Wait()
	unRegister(p.name)

	return
}

// 休眠工作
func (p *Pool) Sleep(seconds int64) bool {
	return p.worker.sleepWorker(seconds)
}

// 添加func到处理队列
func (p *Pool) PushJobFunc(f JobFunc, args ...interface{}) {
	p.worker.job <- &Job{
		f:    f,
		args: args,
	}
}

// 添加任务到处理队列
func (p *Pool) PushJob(job *Job) {
	p.worker.job <- job
}
