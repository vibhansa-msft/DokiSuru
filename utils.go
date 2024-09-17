package main

import (
	"crypto/md5"
	"encoding/base64"
)

func GetBlockID(blockId uint16, md5sum []byte) string {
	// Create a slice which holds blockid followed by md5sum and then create a base64 encoded string
	// This is used to uniquely identify a block

	// Convert the blockId to a slice
	blockIdBytes := make([]byte, 2)
	blockIdBytes[0] = byte(blockId >> 8)
	blockIdBytes[1] = byte(blockId & 0xff)

	// Concat the blockid with the md5sum
	md5sum = append(blockIdBytes, md5sum...)

	// Convert this slice to a base64 encoded string
	return base64.StdEncoding.EncodeToString(md5sum)
}

func ComputeMd5Sum(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}
