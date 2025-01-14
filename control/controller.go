package control

import (
	"FuckEventvwr/config"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

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
			// files = append(files, fmt.Sprint(filepath.Join(Cfg.Path, f.Name())))
			files = append(files, fmt.Sprint(filepath.Join(Cfg.Path, f.Name())))
		}
	}

	var wg sync.WaitGroup
	threads := runtime.NumCPU() * 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go readWork(&wg)
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go writeWork(&wg)
	}
	wg.Wait()
}
