package main

import (
	"FuckEventvwr/cmd"
	"FuckEventvwr/control"
	_ "FuckEventvwr/output"
	"fmt"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	extraArgs := cmd.InitModule(os.Args)
	control.Run(extraArgs)

	fmt.Printf("[SUCCESS] 程序运行时长: %s 任务已完成\n", time.Since(startTime))
}
