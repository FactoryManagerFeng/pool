package pool

import (
	"fmt"
	"strings"
	"sync"
)

var (
	poolMap = make(map[string]*Pool)
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
func NewPool(name string, capacity, running int, worker *Worker) *Pool {
	if p, ok := poolMap[strings.ToLower(name)]; ok {
		return p
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
		worker:   worker,
	}
	return pool
}

func register(p *Pool) {
	poolMap[strings.ToLower(p.name)] = p
}

func unRegister(name string) {
	delete(poolMap, strings.ToLower(name))
}

func Get(name string) *Pool {
	if p, ok := poolMap[strings.ToLower(name)]; ok {
		return p
	}
	return nil
}

// 添加工作线程数
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

// 减少工作线程数
func (p *Pool) DecWorker(num int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if num > p.running {
		p.worker.stopAllWorker()
		p.wg.Wait()
		unRegister(p.name)
		return
	}
	for i := 0; i < num; i++ {
		p.running--
		p.worker.deleteWorker()
	}
	fmt.Println(p)
	return
}

// 开始工作
func (p *Pool) Start() {
	defer p.lock.Unlock()
	p.lock.Lock()
	if p.running > p.capacity {
		return
	}
	var fu = func() { p.wg.Done() }
	for i := 0; i < p.running; i++ {
		p.wg.Add(1)
		p.worker.createWorker(fu)
		fmt.Println("running", p.running)
	}
	register(p)
	return
}

// 停止工作
func (p *Pool) Stop() {
	defer p.lock.Unlock()
	p.lock.Lock()

	p.worker.stopAllWorker()
	p.wg.Wait()
	unRegister(p.name)

	return
}
