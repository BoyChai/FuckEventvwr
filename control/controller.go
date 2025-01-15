package control

import (
	"FuckEventvwr/config"
	"FuckEventvwr/output"
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
	output.Output.Close()
}
