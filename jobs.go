package main

// Job is a job which processes a block
type Job struct {
	Path       string
	Offset     int64
	BlockIndex uint16
	BlockId    string
	Data       []byte
	Md5Sum     []byte
}

// JobHandler interface which needs to be implemented by the job
type JobHandler interface {
	Process(workerId int, job *Job) error
	Start() error
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

func (bh *BaseHandler) GetNext() JobHandler {
	return bh.Next
}

func (bh *BaseHandler) SetNext(next JobHandler) {
	bh.Next = next
}
