package main

// Job interface which needs to be implemented by the job
type Job interface {
	Process(workerId int)
}

// BlockJob is a job which processes a block
type BlockJob struct {
	Path       string
	Offset     int64
	BlockIndex uint16
	BlockId    string
	Data       []byte
	Md5Sum     []byte
}
