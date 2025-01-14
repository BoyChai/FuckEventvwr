package control

import (
	"FuckEventvwr/output"
	"FuckEventvwr/velocidex/evtx"
	"fmt"
	"os"
	"sync"
)

// 锁
var filesLock sync.Mutex
var recordsLock sync.Mutex

// 所有文件
var files []string

// 所有事件
var records []*evtx.EventRecord

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

// 读取线程
func readWork(wg *sync.WaitGroup) {
	for {
		// 读文件
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
		// 解析文件中的块
		chunks, err := evtx.GetChunks(file)
		if err != nil {
			addError("[ERROR] 获取 " + f + " 文件 Chunk 出现错误, 错误信息: " + err.Error())
			continue
		}
		// 解析块中的记录
		for i, c := range chunks {
			if c == nil {
				break
			}
			records, err := c.Parse(0)
			if err != nil {
				addError("[ERROR] 解析" + file.Name() + "文件的第 " + fmt.Sprint(i) + " 个 Chunk 出现错误, 错误信息: " + err.Error())
				continue
			}
			addRecord(records)
		}
		file.Close()
	}
}

// 写入线程
func writeWork(wg *sync.WaitGroup) {
	for {
		record := TakeRecord()
		if record == nil {
			wg.Done()
			return
		}
		output.Output.WriteRecord(record)
	}
}

// 错误处理线程
func errorWork(wg *sync.WaitGroup) {
	for {
		e := <-errorChan
		fmt.Println(e)
		output.Output.WriteError(e)
	}
}
