# 目的
# 协程虽然很好用，但无限制的使用协程会导致系统性能，内存出现问题，该package主要目的是用来控制协程数量


# 方案
# 能设置启动协程数量，能动态增加，减少协程数，能随时停止程序
# Job       任务定义
# Worker    处理任务
# Pool      控制Worker的数量


# 初始化线程池
worker := pool.NewWorker()
newPool := pool.NewPool(name, 100, 1, worker)

# 启动线程池
newPool.Start()

# 向线程池推送数据
worker.PushJobFunc(func(args ...interface{}) pool.State {
    arg := args[0].([]interface{})
    fmt.Println("args", args, "arg", arg)
    time.Sleep(time.Second)
    if arg[0].(int) > 100 {
        return pool.StateOk
    }
    return pool.StateErr
},1)

#获取现有线程池
p := pool.Get(name)
if p == nil {
    fmt.Println("pool is empty", p)
    return
}

#停止线程池消费
p.Stop()

#增加消费线程
p.IncWorker(10)

#减少消费线程
p.DecWorker(10)

#休眠10s
p.Sleep(10)