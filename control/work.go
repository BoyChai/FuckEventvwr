package control

import (
	"FuckEventvwr/velocidex/evtx"
	"sync"
)

// 锁
var filesLock sync.Mutex
var chunksLock sync.Mutex
var recordsLock sync.Mutex

// 所有文件
var files []string

// 所有chunk
var chunks []*evtx.Chunk

// 所有事件
var records []*evtx.EventRecord

func takeFile() string {
	filesLock.Lock()
	defer filesLock.Unlock()

	if len(files) == 0 {
		return ""
	}
	file := files[0]
	files = files[1:]
	return file
}

func takeChunk() *evtx.Chunk {
	chunksLock.Lock()
	defer chunksLock.Unlock()

	if len(chunks) == 0 {
		return nil
	}
	chunk := chunks[0]
	chunks = chunks[1:]
	return chunk
}

func takeRecord() *evtx.EventRecord {
	recordsLock.Lock()
	defer recordsLock.Unlock()

	if len(records) == 0 {
		return nil
	}
	record := records[0]
	records = records[1:]
	return record
}

func fileWork() {

}

func chunkWork() {

}

func recordWork() {

}

func writeWork() {

}
