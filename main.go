package main

import (
	"FuckEventvwr/control"
	"flag"
	"fmt"
)

func main() {

	t := flag.String("t", "p", "指定类型,默认为p,即指定路径模式,s为系统日志")
	p := flag.String("p", "./", "指定路径,默认为当前路径,匹配方式是这个路径下的全部evtx文件")
	o := flag.String("o", "./event.xlsx", "输出位置,默认为当前目录的 event.xlsx")
	c := flag.Bool("c", false, "输出模式,默认为false 追加")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("传入参数有误，请查看使用说明")
		return
	}
	switch *t {
	case "p":
		control.Cfg.Path = *p
	case "s":
		control.Cfg.Path = "C:\\Windows\\System32\\winevt\\Logs"
	default:
		fmt.Println("指定类型错误,请查看帮助")
		return
	}
	control.Cfg.Output = *o
	control.Cfg.Cover = *c

	control.Run()
}
