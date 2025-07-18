package control

import (
	"FuckEventvwr/config"
	"FuckEventvwr/output"
	"FuckEventvwr/velocidex/evtx"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// 活动线程池锁
var threadPoolLock sync.Mutex

// 活动线程池数量
var threadPool int = 0

func Run() {
	Cfg := config.Cfg
	dfs, err := os.ReadDir(Cfg.Path)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, f := range dfs {
		// 拿文件
		if filepath.Ext(f.Name()) == ".evtx" {
			files = append(files, fmt.Sprint(filepath.Join(Cfg.Path, f.Name())))
		}
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
	fmt.Println("日志总量：", count)

	var wg sync.WaitGroup
	threads := runtime.NumCPU()
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go readWork(&wg)
		threadPool++
	}
	for i := 0; i < threads*2; i++ {
		wg.Add(1)
		go writeWork(&wg)
	}
	wg.Wait()
	outCount, err := output.Module.Close()

	if err != nil {
		fmt.Println("日志数量统计失败", err)
		os.Exit(1)
	}
	fmt.Println("模块收集日志数据量：", outCount)
}
