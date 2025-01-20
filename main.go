package main

import (
	"FuckEventvwr/config"
	"FuckEventvwr/control"
	"FuckEventvwr/output"
	"flag"
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	outputFileName := fmt.Sprintf("%sevent-%s.xlsx", "./", time.Now().Format("2006-01-02-15-04-05"))

	// 解析命令行参数
	t := flag.String("t", "p", "指定类型,默认为p,即指定路径模式,s为系统日志")
	p := flag.String("p", "./", "指定路径,默认为当前路径,匹配方式是这个路径下的全部evtx文件")
	o := flag.String("o", outputFileName, "输出位置,默认为当前目录的 event-年-月-日-时-分-秒.xlsx")
	m := flag.Int("m", 0, "输出模式,默认0,0为只打印处理好的数据,1则处理好的数据+原始数据都放到一个xlsx中,2则只输出原始数据")
	eu := flag.String("eu", "", "ESURL,设置ESURL则只往ES中打入数据,默认为空")

	flag.Parse()

	Cfg := config.Cfg
	switch *t {
	case "p":
		Cfg.Path = *p
	case "s":
		Cfg.Path = "C:\\Windows\\System32\\winevt\\Logs"
	default:
		fmt.Println("指定类型错误,请查看帮助")
		return
	}

	Cfg.Output = *o
	Cfg.Mode = *m
	Cfg.EsURL = *eu

	output.InitOutput()
	control.Run()

	fmt.Printf("[SUCCESS] 程序运行时长: %s 任务已完成\n", time.Since(startTime))
}
