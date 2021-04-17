package pool

type JobFunc func(args ...interface{}) JobResult

type Job struct {
	f    JobFunc
	args []interface{}
}

func (job *Job) execute() JobResult {
	return job.f(job.args)
}

func NewJob(f JobFunc, args ...interface{}) *Job {
	return &Job{
		f:    f,
		args: args,
	}
}

type JobResult struct {
	State State
	Err   error
}
