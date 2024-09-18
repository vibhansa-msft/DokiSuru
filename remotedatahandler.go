package main

type RemoteDataHandler struct {
}

func NewRemoteDataHandler() *RemoteDataHandler {
	return &RemoteDataHandler{}
}

// Process the block
func (ldh *RemoteDataHandler) Process(workerId int, bj *Job) error {
	return nil
}

func (ldh *RemoteDataHandler) Start(schedule func(job *Job)) {

}
