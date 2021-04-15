package pool

type JobFunc func(args ...interface{}) State

type Job struct {
	f    JobFunc
	args []interface{}
}

func (job *Job) execute() State {
	return job.f(job.args)
}

func (job *Job) Execute() State {
	return job.f(job.args)
}

func NewJob(f JobFunc, args ...interface{}) *Job {
	return &Job{
		f:    f,
		args: args,
	}
}
