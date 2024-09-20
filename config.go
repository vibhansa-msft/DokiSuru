package main

type dokiSuruConfig struct {
	WorkerCount int
	BlockSize   uint64
	Path        string
	Validate    bool
}

var config dokiSuruConfig
