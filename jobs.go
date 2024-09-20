package main

// Job is a job which processes a block
type Job struct {
	BlockIndex uint16
	NoOfBlocks uint16
	BlockId    string
	Data       []byte
	Md5Sum     []byte
}

// JobHandler interface which needs to be implemented by the job
type JobHandler interface {
	Start() error
	Stop() error

	Process(workerId int, job *Job) error

	Enqueue(job *Job)

	GetNext() JobHandler
	SetNext(JobHandler)
}

type BaseHandler struct {
	Next   JobHandler
	Worker *WorkerPool
}

func (bh *BaseHandler) Process(workerId int, job *Job) error {
	return nil
}

func (bh *BaseHandler) Start() error {
	return nil
}

func (bh *BaseHandler) Stop() error {
	return nil
}

func (bh *BaseHandler) GetNext() JobHandler {
	return bh.Next
}

func (bh *BaseHandler) SetNext(next JobHandler) {
	bh.Next = next
}

func (bh *BaseHandler) Enqueue(job *Job) {
	bh.Worker.AddJob(job)
}
