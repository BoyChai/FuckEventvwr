package control

import (
	"FuckEventvwr/output"
	"FuckEventvwr/velocidex/evtx"
	"os"
	"sync"
)

// 锁
var filesLock sync.Mutex
var chunksLock sync.Mutex
var recordsLock sync.Mutex

// 所有文件
var files []string

// 所有 chunk
var chunks []*evtx.Chunk

// 所有事件
var records []*evtx.EventRecord

// 添加文件
func addFile(file string) {
	filesLock.Lock()
	defer filesLock.Unlock()
	files = append(files, file)
}

// 取出文件
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

// 添加 chunk
func addChunk(chunk []*evtx.Chunk) {
	chunksLock.Lock()
	defer chunksLock.Unlock()
	for _, c := range chunk {
		chunks = append(chunks, c)
	}
}

// 取出 chunk
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

// 添加记录
func addRecord(record []*evtx.EventRecord) {
	recordsLock.Lock()
	defer recordsLock.Unlock()
	for _, r := range record {
		records = append(records, r)
	}
}

// 取出记录
func TakeRecord() *evtx.EventRecord {
	recordsLock.Lock()
	defer recordsLock.Unlock()

	if len(records) == 0 {
		return nil
	}
	record := records[0]
	records = records[1:]
	return record
}

// func fileWork(fc chan string, wg *sync.WaitGroup) {
func fileWork(wg *sync.WaitGroup) {

}

// 读取文件拿chunk
func chunkWork(wg *sync.WaitGroup) {
	for {
		f := takeFile()
		if f == "" {
			wg.Done()
			return
		}
		file, err := os.Open(f)
		if err != nil {
			addError("[ERROR] 打开" + f + "文件失败:" + err.Error())
			continue
		}
		// 获取文件中的所有块（Chunk）
		chunks, err := evtx.GetChunks(file)
		file.Close()
		if err != nil {
			addError("[ERROR] 获取 " + f + " 文件 Chunk 出现错误, 错误信息: " + err.Error())
			continue
		}
		addChunk(chunks)
	}
}

// 读chunk拿记录
func recordWork(wg *sync.WaitGroup) {
	for {
		c := takeChunk()
		if c == nil {
			wg.Done()
			return
		}
		records, err := c.Parse(0)
		if err != nil {
			addError("[ERROR] 解析 Chunk 出现错误, 错误信息: " + err.Error())
			continue
		}
		addRecord(records)
	}

}

func writeWork(wg *sync.WaitGroup) {
	for {
		r := TakeRecord()
		if r == nil {
			wg.Done()
			return
		}
		err := output.Output.Write(r)
		if err != nil {
			addError("[ERROR] 写入数据出现错误, 错误信息: " + err.Error())
			continue
		}
	}
}
