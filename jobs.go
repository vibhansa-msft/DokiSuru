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

// Job interface which needs to be implemented by the job
type JobHandler interface {
	Process(workerId int, job *Job) error
}
