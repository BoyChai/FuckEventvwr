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

// 所有文件
var files []string

// 事件通道
var recordChan chan *evtx.EventRecord = make(chan *evtx.EventRecord, 10000)

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

// 读取线程
func readWork(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// 读文件
		f := takeFile()
		if f == "" {
			threadPoolLock.Lock()
			threadPool--
			threadPoolLock.Unlock()
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
			for _, record := range records {
				record.FileName = file.Name()
				recordChan <- record
			}
		}
		file.Close()
	}
}

// 写入线程
func writeWork(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// 有数据
		case record := <-recordChan:
			err := output.Module.WriteRecord(record)
			if err != nil {
				addError("[ERROR] 写入记录失败, 错误信息: " + err.Error())
			}
		// 没数据
		default:
			if threadPool == 0 {
				return
			}
		}
	}
}

// 数据校验
func checkData(files []string) (bool, int, error) {
	outCount, err := output.Module.Count()
	if err != nil {
		return false, 0, err
	}
	var count int
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		c, err := evtx.CountLogs(f)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		count += c
	}
	if outCount != count {
		return false, outCount, nil
	}
	return true, outCount, nil
}
